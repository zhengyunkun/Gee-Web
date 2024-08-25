# A Simple Stupid Server in golang

## with gee web frame

### 定义请求处理函数

```go
type HandlerFunc func(http.ResponseWriter, *http.Request)
```

### 在Web框架中，router用于存储不同的URL路径(string)与处理这些路径的函数(HandlerFunc)之间的映射关系

```go
type Engine struct {
    router map[string]HandlerFunc
    // map string to HandlerFunc
}
```

### 将Engine封装成http.Handler

#### 实现ServeHTTP方法

```go
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    key := req.Method + "-" + req.URL.Path
    if handler, ok := engine.router[key]; ok {
        handler(w, req)
    } else {
        fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
    }
}
```

### 为Engine实现一系列构造函数和处理函数

```go
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
    key := method + "-" + pattern
    engine.router[key] = handler
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
    engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
    engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) RUN(addr string) (err error) {
    return http.ListenAndServe(addr, engine)
}
```
