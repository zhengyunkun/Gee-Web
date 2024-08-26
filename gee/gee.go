package gee

import (
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type (
	RouterGroup struct {
		prefix	    string
		// middlewares []HandlerFunc // support middleware
		parent      *RouterGroup  // support nesting
		engine      *Engine       // all groups share a Engine instance
	}

	Engine struct {
		*RouterGroup
		router 		*router
		groups 		[]*RouterGroup // store all groups
	}
)

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}					// create a new Engine object
	engine.RouterGroup = &RouterGroup{engine: engine}		// create a new RouterGroup object
	engine.groups = []*RouterGroup{engine.RouterGroup}		// add the default group to groups
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup {
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// addRoute moves to RouterGroup
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
// pattern is the url path like "/" or "/hello"
// HandlerFunc is the function to handle the request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
// addr is the server address like ":8080"
// the (err error) is the return value of ListenAndServe
func (engine *Engine) RUN(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP implement all the logic of a http request
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}