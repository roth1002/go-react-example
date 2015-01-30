## go-react-example
Example project to show how to render React app on Golang server-side.  
It's very similar project as [andreypopp/react-quickstart](https://github.com/andreypopp/react-quickstart)(please see it first), but with Go and some key benefits.  

## What it contains?
- [go-duktape]() bindings for a thin, embeddable javascript engine;
- [gin]() framework;
- [staticbin]() middleware for gin, to serve embedded static files;
- [config]() package, to define config and parse flags and environment variables
- [optional] live code reloading, by `fswatch`, avilable for OSX and linux.


## Benefits & differences
First of all and most important, rendering is **fully synchronous**. There is no need to use react-async. Because on the server side the code executes in synchronous mode. This is duktape specific. Also, there is a binding between superagent javascript package and golang server side function. Often it is called monkey patching. That means that you do http request and it handles by golang function. As a consequence, there is no need to do a http request from the server to the server. Now it it works between you react application and server side application directly. And possible to handle requests with user session as well.

Also this project allows you embed all static files. So, you have one executable file of you application. Cross compiling is also avilable with this approach, but this is not the point.


## Install

```
$ # clone the repo
$ git clone https://github.com/olebedev/go-react-example && cd $_
$ # fetch dependencies
$ go get ./...
$ # install go-bindata to embed static files
$ go get -u github.com/jteeuwen/go-bindata/...
```

Now you ready to start.

```
$ make
$ go run *.go
```

> If you have `fswatch`, you can type this `make serve` and project will be reload every time when you change any `*.go` and `static/*` files.  
> `fswatch` avilable for OSX and linux.

## Benchmarks
