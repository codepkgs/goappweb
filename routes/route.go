package routes

import "github.com/gin-gonic/gin"

func Init() *gin.Engine {
	r := gin.New()
	return r
}

func RegisterMiddleware(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	r.Use(middlewares...)
}

type Registerer interface {
	Register(engine *gin.Engine) gin.HandlerFunc
}

func RegisterRoutes(r *gin.Engine, routes []Registerer) {
	for _, route := range routes {
		route.Register(r)
	}
}
