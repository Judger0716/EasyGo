package main

import (
	"EasyCache"
	"EasyGin"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *EasyCache.Group {
	return EasyCache.NewGroup("scores", 2<<10, EasyCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *EasyCache.Group) {
	peers := EasyCache.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("EasyCache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *EasyCache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

		}))
	log.Println("frontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {

	r := EasyGin.Default()
	r.Use(EasyGin.Logger()) // global midlleware
	r.LoadHTMLGlob("EasyGin/templates/*")
	r.Static("/assets", "./static")

	// frontend
	r.GET("/", func(c *EasyGin.Context) {
		c.HTML(http.StatusOK, "css.tmpl", "<h1>Hello EasyGin</h1>")
	})

	// distributed nodes
	cacheGroup := createGroup()
	for i := 0; i < 3; i++ {
		node := EasyGin.Default()
		node.Use(EasyGin.Logger())
		node.GET("/api", func(c *EasyGin.Context) {
			key := c.Query("key")
			log.Println(key)
			view, err := cacheGroup.Get(key)
			if err != nil {
				c.Fail(500, "Internel Server Error")
				return
			}
			c.Data(200, view.ByteSlice())
			// w.Header().Set("Content-Type", "application/octet-stream")
			// w.Write(view.ByteSlice())
		})
		go node.Run(":" + strconv.Itoa(8001+i))
	}

	// api server
	cache := r.Group("/cache")
	cache.Use(EasyGin.CacheLogger()) // v2 group middleware
	{
		cache.GET("/hello/:name", func(c *EasyGin.Context) {
			// expect /hello/EasyGinktutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		cache.GET("/api", func(c *EasyGin.Context) {
			key := c.Query("key")
			log.Println(key)
			view, err := cacheGroup.Get(key)
			if err != nil {
				c.Fail(500, "Internel Server Error")
				return
			}
			c.Data(200, view.ByteSlice())
			// w.Header().Set("Content-Type", "application/octet-stream")
			// w.Write(view.ByteSlice())
		})
	}

	r.Run(":9999")
}
