package main

import (
	"net/http"

	"easygin"
)

func main() {
	r := easygin.Default()
	r.Use(easygin.Logger()) // global midlleware
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")
	r.GET("/", func(c *easygin.Context) {
		c.HTML(http.StatusOK, "css.tmpl", "<h1>Hello easygin</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(easygin.OnlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *easygin.Context) {
			// expect /hello/easyginktutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
