# A Simple Stupid Server in golang

## with gee web frame

### Trie树实现动态路由

给定

```go
// 给定现有路由
routes := []string{
    "/p/:lang/doc",
    "/p/:lang/tutorial",
    "/p/:lang/intro",
    "/static/*filepath",
}

// 查找路由
searchRoutes := []string{
    "/p/go/doc",
    "/p/python/tutorial",
    "/static/css/style.css",
}
```
