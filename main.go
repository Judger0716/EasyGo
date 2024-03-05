package main

import (
	"EasyCache"
	"EasyGin"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "EasyCache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	r := EasyGin.Default()
	// r.Use(EasyGin.Logger()) // global midlleware

	// global config
	apiAddr := "http://localhost:9999"
	addrs := []string{
		"http://localhost:8001",
		"http://localhost:8002",
		"http://localhost:8003",
	}
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	cacheGroup := EasyCache.CreateGroup()
	if api {

		// api server
		cache := EasyGin.Default()
		// cache.Use(EasyGin.CacheLogger())

		// load resources
		cache.LoadHTMLGlob("EasyGin/templates/*")
		cache.Static("/assets", "./static")

		// frontend
		cache.GET("/", func(c *EasyGin.Context) {
			c.HTML(http.StatusOK, "css.tmpl", "<h1>Hello EasyGin</h1>")
		})

		cache.GET("/hello/:name", func(c *EasyGin.Context) {
			// expect /hello/EasyGinktutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		// 返回客户端
		cache.GET("/api", func(c *EasyGin.Context) {
			key := c.Query("key")
			view, err := cacheGroup.Get(key)
			if err != nil {
				c.Fail(500, "Internel Server Error")
				return
			}
			c.Data(200, view.ByteSlice())
			// w.Header().Set("Content-Type", "application/octet-stream")
			// w.Write(view.ByteSlice())
		})

		go cache.Run(":9999")
		log.Printf("Fontend server is running at %s\n", apiAddr)

	}

	// register group
	peers := EasyCache.NewHTTPPool(addrMap[port])
	peers.Set(addrs...)
	cacheGroup.RegisterPeers(peers)
	log.Printf("EasyCache is running at %s\n", addrMap[port])

	r.GET("/_easycache/scores/:name", func(c *EasyGin.Context) {

		// 未通过Key=Tom传输，模糊匹配并解析
		// 原ServeHTTP
		// log.Println(c.Req.URL.Path)
		parts := strings.SplitN(c.Req.URL.Path[len(EasyCache.DefaultBasePath):], "/", 2)
		if len(parts) != 2 {
			c.Fail(http.StatusBadRequest, "Parse param failed.")
			return
		}

		groupName := parts[0]
		key := parts[1]

		group := EasyCache.GetGroup(groupName)
		if group == nil {
			c.Fail(http.StatusNotFound, fmt.Sprintf("No such group: %s", groupName))
			return
		}

		view, err := group.Get(key)
		if err != nil {
			c.Fail(500, "Internel Server Error")
			return
		}
		c.Data(http.StatusOK, view.ByteSlice())
		// w.Header().Set("Content-Type", "application/octet-stream")
		// w.Write(view.ByteSlice())
	})
	r.Run(":" + strconv.Itoa(port))
}
