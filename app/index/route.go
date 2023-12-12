package index

import (
	"github.com/gin-gonic/gin"
)

func Routes(rg *gin.RouterGroup) {
	index := rg.Group("/index")
	index.GET("/", HomeHandler)
}
