package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "hello gin")
}
