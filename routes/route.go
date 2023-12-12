package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterMiddlewares(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	r.Use(middlewares...)
}

func Init() *gin.Engine {
	r := gin.New()

	return r
}

type Register func(r *gin.RouterGroup)

func RegisterRoutes(register Register, rg *gin.RouterGroup) {
	register(rg)
}
