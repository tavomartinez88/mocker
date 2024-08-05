package main

import (
	"encoding/json"
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

	// Ruta para agregar nuevos endpoints
	r.POST("/add-endpoint", func(c *gin.Context) {
		var routeRequest struct {
			Method  string `json:"method"`
			Pattern string `json:"pattern"`
			Body    string `json:"body"`
			Status  int    `json:"status"`
		}

		if err := c.ShouldBindJSON(&routeRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var response map[string]interface{}
		if err := json.Unmarshal([]byte(routeRequest.Body), &response); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al decodificar JSON"})
			return
		}

		handler := func(c *gin.Context) {
			c.JSON(routeRequest.Status, response)
		}

		dynamicRoutesMu.Lock()
		dynamicRoutes[routeRequest.Pattern] = dynamicRoute{
			method:  routeRequest.Method,
			pattern: routeRequest.Pattern,
			handler: handler,
		}
		dynamicRoutesMu.Unlock()

		c.JSON(http.StatusOK, gin.H{"message": "Endpoint added"})
	})

	// Ruta para eliminar endpoints
	r.POST("/remove-endpoint", func(c *gin.Context) {
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
	})

	handlerDynamicRoutes(r)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
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
