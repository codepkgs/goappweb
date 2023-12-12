package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func indexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "hello gin")
	}
}
