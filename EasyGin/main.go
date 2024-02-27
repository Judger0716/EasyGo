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

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(ctx *easygin.Context) {
			ctx.HTML(http.StatusOK, "<h1>Helle Group1<h1>")
		})

		v1.GET("/hello", func(ctx *easygin.Context) {
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
		})

	}

	r.GET("/hello", func(ctx *easygin.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	r.GET("/hello/:name", func(ctx *easygin.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
	})

	r.GET("/assets/*filepath", func(ctx *easygin.Context) {
		ctx.JSON(http.StatusOK, easygin.H{"filepath": ctx.Param("filepath")})
	})

	r.POST("/login", func(ctx *easygin.Context) {
		ctx.JSON(http.StatusOK, easygin.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	r.Run(":9999")
}
