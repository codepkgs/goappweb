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

func RegisterRoutes(rg *gin.RouterGroup, registers ...Register) {
	for _, register := range registers {
		register(rg)
	}
}
