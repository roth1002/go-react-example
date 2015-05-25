package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olebedev/go-duktape"
)

type duktapePool struct {
	ch chan *duktape.Context
}

func (o *duktapePool) get() *duktape.Context {
	return <-o.ch
}

func (o *duktapePool) put(ot *duktape.Context) {
	o.ch <- ot
}

func newDukPool(size int, engine *gin.Engine) *duktapePool {
	pool := &duktapePool{
		ch: make(chan *duktape.Context, size),
	}
loop:
	for {
		select {
		case pool.ch <- newDukContext(engine):
		default:
			break loop
		}

	}
	return pool
}

// Loads bundle.js to context
func newDukContext(engine *gin.Engine) *duktape.Context {
	vm := duktape.Default()
	if err := vm.PevalString(`var self = {}, console = {log:print,warn:print,error:print,info:print}`); err != nil {
		panic(err.(*duktape.Error).Message)
	}
	app, err := Asset("bundle.js")
	panicIf(err)
	if err := vm.PevalString(string(app)); err != nil {
		panic(err.(*duktape.Error).Message)
	}
	panicIf(vm.PushGlobalGoFunction("__request__", request(engine)))
	if err := vm.PevalString(superagentBinding); err != nil {
		panic(err.(*duktape.Error).Message)
	}

	// test case for superagent
	// vm.EvalString(`
	// print("start ============================================");
	// self.superagent.get('/api/v1/config', function(err, response){
	//   console.log("response", err, JSON.stringify(response.body, null, 2));
	// });
	// print("stop =============================================");
	// `)

	return vm
}

type __react__ struct {
	pool   *duktapePool
	engine *gin.Engine
}

func (r *__react__) init() {
	r.pool = newDukPool(conf.UInt("duktape.poolSize"), r.engine)
}

func (r *__react__) handle(c *gin.Context) {
	var v string
	vm := r.pool.get()
	vm.PushGlobalObject()
	vm.PevalString(`self.React.renderToString(self.App({path:'` + c.Request.URL.Path + `'}));`)
	v = vm.SafeToString(-1)
	vm.Pop()
	r.pool.put(vm)

	c.Writer.WriteHeader(200)
	c.Writer.Header().Add("Content-Type", "text/html")
	c.Writer.Write([]byte("<!doctype html>\n" + v))
}

func (r *__react__) bind(ro *gin.Engine) {
	r.engine = ro
	r.init()
	ro.Use(r.handle)
}

var react = __react__{}

func debug(st string, values ...interface{}) {
	log.Printf(st, values...)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func pushToCtxAndExit(ctx *duktape.Context, res map[string]interface{}) int {
	b, _ := json.Marshal(res)
	ctx.PushString(string(b))
	return 1
}

func request(engine *gin.Engine) func(*duktape.Context) int {
	return func(ctx *duktape.Context) int {
		res := map[string]interface{}{}
		req := map[string]interface{}{}
		// read first argument
		err := json.Unmarshal([]byte(ctx.GetString(1)), &req)
		if err != nil {
			res["code"] = http.StatusInternalServerError
			res["body"] = err.Error()
			return pushToCtxAndExit(ctx, res)
		}

		if ctx.GetTop() != 2 {
			res["code"] = http.StatusInternalServerError
			res["body"] = "request object not found"
			return pushToCtxAndExit(ctx, res)
		}

		method := req["method"].(string)
		url := req["url"].(string)
		// it needs to make crossdomain request from server side, if need
		if strings.HasPrefix(url, "//") || strings.HasPrefix(url, "http") {
			url = "/away?url=" + url
		}

		header := make(map[string]interface{})
		if h, ok := req["headers"]; ok && h != nil {
			header = h.(map[string]interface{})
		}
		var body io.Reader
		if b, ok := req["body"]; ok {
			if _b, ok := b.(string); ok {
				body = bytes.NewReader([]byte(_b))
			}
		}

		response := httptest.NewRecorder()
		request, err := http.NewRequest(
			method,
			url,
			body,
		)

		// TODO: copy cookie to request
		for k, v := range header {
			request.Header.Add(k, v.(string))
		}
		request.Header.Add("X-Server-React", "true")

		// make request
		engine.ServeHTTP(response, request)

		rheaders := map[string][]string{}
		for k, v := range response.HeaderMap {
			rheaders[k] = v
		}
		res["headers"] = rheaders
		res["code"] = response.Code
		res["body"] = response.Body.String()
		if err != nil {
			res["code"] = http.StatusInternalServerError
			res["body"] = err.Error()
		}

		return pushToCtxAndExit(ctx, res)
	}
}

// Embed javascript to make project folder simpler
const superagentBinding = `
/**
  * Mostly copied from original superagent source
  */
var request = self.superagent.Request;
// define func to avoid exeption
function isHost(obj) {
  var str = {}.toString.call(obj);

  switch (str) {
    case '[object File]':
    case '[object Blob]':
    case '[object FormData]':
      return true;
    default:
      return false;
  }
}

// Overwrite Request end method
request.prototype.end = function(fn) {
  var query = this._query.join('&');
  var data = this._formData || this._data;

  // store callback
  this._callback = fn || function(){};

  // querystring
  if (query) {
    query = request.serializeObject(query);
    this.url += ~this.url.indexOf('?')
      ? '&' + query
      : '?' + query;
  }

  // body
  if ('GET' != this.method && 'HEAD' != this.method && 'string' != typeof data && !isHost(data)) {
    // serialize stuff
    var serialize = request.serialize[this.getHeader('Content-Type')];
    if (serialize) data = serialize(data);
  }

  var headers = {};
  // set header fields
  for (var field in this.header) {
    if (null == this.header[field]) continue;
    headers[field] = this.header[field];
  }

  this.emit('request', this);
  this.__response__ = JSON.parse(__request__(JSON.stringify({
    url: this.url,
    method: this.method,
    headers: headers,
    body: data
  })))
  // Generate a response & call back
  var err = null;
  var res = null;

  try {
    res = new Response(this);
  } catch(e) {
    err = new Error('Parser is unable to parse the response');
    err.parse = true;
    err.original = e;
  }

  if (res) {
    this.emit('response', res);
  }
  this._callback(err, res)

  this._callback = function(){};
  try {
    this.emit('end');
  } catch (e) {};
  return this;
};

Response.prototype = Object.create(self.superagent.Response.prototype);
for (var key in self.superagent.Response.prototype) {
  Response.prototype[key] = self.superagent.Response.prototype[key];
};

// Overwrite the Response constructor
function Response(req, options) {
  options = options || {};
  this.req = req;
  this.__response__ = this.req.__response__;
  this.text = this.req.method !='HEAD'
     ? this.__response__.body
     : null;

  this.setStatusProperties(this.__response__.code);

  this.header = this.headers = this.__response__.headers;
  for (var key in this.header) {
    this.header[key.toLowerCase()] = this.header[key].join();
    delete this.header[key]
  }
  this.setHeaderProperties(this.header);

  this.body = this.req.method != 'HEAD'
    ? this.parseBody(this.text)
    : null;
};
`
