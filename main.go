package main

import (
	"fmt"
	"net/http"
	"gee"
)

func main() {
	r := gee.New()

	// Add a route for GET request
	// ResponseWriter is used to send data to the client
	// Request is used to get data from the client
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	// Add a route for GET request
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.RUN(":8080")
}