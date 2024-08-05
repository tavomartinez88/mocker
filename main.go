package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type dynamicRoute struct {
	method  string
	pattern string
	handler gin.HandlerFunc
}

var (
	dynamicRoutes   = make(map[string]dynamicRoute)
	dynamicRoutesMu sync.RWMutex
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
		c.JSON(http.StatusOK, "healthy")
	})

	r.POST("/add-endpoint", func(c *gin.Context) {
		addEndpoint(c)
	})

	r.POST("/remove-endpoint", func(c *gin.Context) {
		removeEndpoint(c)
	})

	handlerDynamicRoutes(r)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}

func removeEndpoint(c *gin.Context) {
	var route struct {
		Pattern string `json:"pattern"`
	}
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dynamicRoutesMu.Lock()
	delete(dynamicRoutes, route.Pattern)
	dynamicRoutesMu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Endpoint removed"})
}

func addEndpoint(c *gin.Context) {
	var routeRequest struct {
		Method  string      `json:"method"`
		Pattern string      `json:"pattern"`
		Body    interface{} `json:"body"`
		Status  int         `json:"status"`
	}

	if err := c.ShouldBindJSON(&routeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if routeRequest.Status < 100 || routeRequest.Status > 599 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status code not supported"})
		return
	}

	bodyMap, ok := routeRequest.Body.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body format"})
		return
	}

	handler := func(c *gin.Context) {
		c.JSON(routeRequest.Status, bodyMap)
	}

	dynamicRoutesMu.Lock()
	dynamicRoutes[routeRequest.Pattern] = dynamicRoute{
		method:  routeRequest.Method,
		pattern: routeRequest.Pattern,
		handler: handler,
	}
	dynamicRoutesMu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Endpoint added"})
}

func handlerDynamicRoutes(r *gin.Engine) gin.IRoutes {
	return r.Use(func(c *gin.Context) {
		dynamicRoutesMu.RLock()
		route, exists := dynamicRoutes[c.Request.URL.Path]
		dynamicRoutesMu.RUnlock()

		if exists && c.Request.Method == route.method {
			route.handler(c)
			c.Abort()
			return
		}

		c.Next()
	})
}
