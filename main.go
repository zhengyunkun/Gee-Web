package main

import (
	"gee"
	"net/http"
	"log"
	"time"
)

func MiddleWareForV2() gee.HandlerFunc {
	return func (c *gee.Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Fail(http.StatusInternalServerError, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	// return a new Engine object
	r.Use(gee.Logger())
	// add global middleware

	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Zhengyunkun's Index Page</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Zhengyunkun!</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=zhengyunkun
			c.String(http.StatusOK, "Hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(MiddleWareForV2())
	// add middleware to v2 group
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "Hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.POST("/login", func(c *gee.Context){
			c.JSON(http.StatusOK, gee.H {
				"username": c.PostForm("username"),
				"sex": c.PostForm("sex"),
				"age": c.PostForm("age"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context){
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"sex": c.PostForm("sex"),
			"age": c.PostForm("age"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.RUN(":8080")
}