## go-react-example
This is an example of project which shows how to render React app on Golang server-side.  
It's very similar to [andreypopp/react-quickstart](https://github.com/andreypopp/react-quickstart)(please see it first) project, but with some Go and other benefits.  

## What it contains?
- [go-duktape](https://github.com/olebedev/go-duktape) bindings for a thin, embeddable javascript engine
- [gin](https://github.com/gin-gonic/gin) framework
- [staticbin](https://github.com/olebedev/staticbin) middleware for gin, to serve embedded static files
- [config](https://github.com/olebedev/config) package, to define config, parse flags and environment variables
- [optional] live code reloading, by `fswatch`, avilable for OSX and linux


## Benefits 
First of all and most important, rendering is **fully synchronous**. There is no need to use react-async. Because on the server side the code executes in synchronous mode. This is duktape specific. Also, there is a binding between superagent javascript package and golang server side function. That means that you do http request and it is processed by golang function. As a consequence, there is no need to do real http request from the server to the same. Now it works between your react application and server side application directly. And it is possible to process requests with user session as well.

Also this project allows you to embed all static files. So, you have one executable file of you application.

## Install

```
$ git clone https://github.com/olebedev/go-react-example && cd $_
$ go get ./...
$ go get -u github.com/jteeuwen/go-bindata/...
$ npm i
```

Now you ready to start.

```
$ make
$ go run *.go
```

> If you have `fswatch`, you can type this `make serve` and project will be reloaded every time when you change any of `*.go` or `static/*` files.  
> `fswatch` is avilable for OSX and linux.
