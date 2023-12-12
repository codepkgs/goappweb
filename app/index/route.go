package index

import "github.com/gin-gonic/gin"

type Routes struct {
	R *gin.Engine
}

func (i *Routes) Register() {
	i.R.GET("/", indexHandler())
}
