package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olebedev/staticbin"
)

func main() {
	// Parse flags and environment variables.
	// See: conf.go file for more details.
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
