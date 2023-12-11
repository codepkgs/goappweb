package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-hezhang/gowebapp/logger"

	"github.com/x-hezhang/gowebapp/settings"
)

// Go Web较通用的脚手架模板

func main() {
	// 配置文件初始化
	if err := settings.Init("config.toml"); err != nil {
		log.Fatalf("init settings failed! %v\n", err)
	} else {
		fmt.Println("init settings success!")
	}

	// 日志初始化
	if err := logger.Init(); err != nil {
		log.Fatalf("init logger failed! %v\n", err)
	} else {
		fmt.Println("init logger success!")
	}

	// MySQL连接
	// Redis连接
	// 注册路由
	// 服务启动（优雅关闭）

	gin.SetMode(settings.Conf.AppConfig.Mode)
	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	r.Run(fmt.Sprintf(":%v", settings.Conf.AppConfig.Port))
}
