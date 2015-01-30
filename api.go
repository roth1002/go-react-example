package main

import (
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

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
