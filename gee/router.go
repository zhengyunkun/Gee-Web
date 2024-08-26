package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots	  map[string]*node  // key is the method, value is the root of the trie tree, 
	// roots['GET'] is the root of the trie tree for GET requests, roots['POST'] is the root of the trie tree for POST requests
	handlers  map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:       make(map[string]*node),
		handlers:    make(map[string]HandlerFunc),
	}
}

// Only one "*" is allowed in the path, return parsed parts of the path
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// Add a route to the router
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := parsePattern(pattern)
	key := method + "-" + pattern //like "GET-/p/:lang/doc"
	_, ok := r.roots[method] 
	// Check if the root of the trie tree for the method exists
	// if not, return nil, otherwise return the root of the trie tree
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0) 	// Add the pattern to the trie tree
	r.handlers[key] = handler					// Add the handler to the handlers map 
}

// Get the route from the router, return values like "GET-/p/:lang/doc", map["lang":"go"]
// input path is like "p/go/doc", in trie there is a node with pattern "/p/:lang/doc"
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0) // the leaf node of the trie tree

	if n != nil {
		parts := parsePattern(n.pattern) // parts is like ["p", ":lang", "doc"] or ["assets", "*filepath"]
		for index, part := range parts { // part range from "p" to ":lang" to "doc"
			if part[0] == ':' {
				// when part is like ":lang", params["lang"] = searchParts[index] = "go"
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 { // part
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}

		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

