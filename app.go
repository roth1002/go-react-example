package main

import (
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/olebedev/config"
	"github.com/olebedev/staticbin"
)

// Base configuration
var conf, _ = config.ParseYaml(`
env: development
app:
  name: go react example
api:
  prefix: /api/v1
duktape:
  poolSize: 5
debug: true
port: 5000
title: Go React Example
`)

type __api__ struct{}

var api = __api__{}

func (api __api__) bind(r *gin.RouterGroup) {
	r.GET("/users/:username", api.username)
	r.GET("/config", api.config)
}

func (_ __api__) config(c *gin.Context) {
	c.JSON(200, conf.Root)
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func (_ __api__) username(c *gin.Context) {
	c.JSON(200, map[string]string{
		"username": c.Params.ByName("username"),
		"name":     upperFirst(c.Params.ByName("username")),
	})
}

func main() {
	// Parse flags and environment variables.
	conf.Flag().Env()

	router := gin.Default()

	// Serve assets from binary
	router.Use(staticbin.Static(Asset, staticbin.Options{
		SkipLogging: true,
		Dir:         "/static",
	}))

	// Attach api, see: api.go.
	api.bind(router.Group(conf.UString("api.prefix")))

	// For all other requests, see: react.go.
	react.bind(router)

	// Start listening
	router.Run(":" + conf.UString("port"))
}
