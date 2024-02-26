package main

import (
	"net/http"

	"easygin"
)

func main() {

	r := easygin.New()
	r.GET("/", func(ctx *easygin.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello EasyGin</h1>")
	})

	r.GET("/hello", func(ctx *easygin.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	r.POST("/login", func(ctx *easygin.Context) {
		ctx.JSON(http.StatusOK, easygin.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	r.Run(":9999")
}
