package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// original objects
	Writer  		http.ResponseWriter
	Req				*http.Request
	// Request Info
	Path			string
	Method			string
	Params			map[string]string
	// Response Info
	StatusCode 		int
	// Middleware
	handlers 		[]HandlerFunc	// store all the middleware
	index 			int				// index of the current middleware
}

// constructor for Context, return a new Context object
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context {
		Writer: 		w,
		Req: 			req,
		Path: 			req.URL.Path,
		Method: 		req.Method,
		index:			-1, // means the middleware has not been executed
	}
}	

// when the middleware is executed, the Next() function is called to execute the next middleware
func (c *Context) Next() {
	c.index ++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c) // pass the context to the middleware and execute it
	}
}

// Fail method is used to handle errors, it writes the error message to the response in JSON format
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers) // skip all the other middlewares
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

// PostForm gets a form value
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Get a key from the URL query string
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Set status code for the response
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// Set a header for the response
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// Send a response with a string
func (c *Context) String(code int, format string, values ...interface{}) {
	// values ...interface{} is a variadic parameter
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...))) // expand values to a list of values
}

// change object to JSON format and send it as a response
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	// automatically write the encoded JSON data to the response
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil { 
		http.Error(c.Writer, err.Error(), 500)
	}
}

// write date to the response in []bytes
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// write HTML to the response
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}