package main

import (
	"fmt"
	"log"

	"github.com/x-hezhang/gowebapp/settings"
)

// Go Web较通用的脚手架模板

func main() {
	// 配置文件初始化
	if err := settings.InitConfig("config.ini"); err != nil {
		log.Fatalf("init settings failed! %v\n", err)
	} else {
		fmt.Println("init settings success!")
	}
	// 配置文件（本地或远程）
	// 日志初始化
	// MySQL连接
	// Redis连接
	// 注册路由
	// 服务启动（优雅关闭）

	fmt.Printf("%#v", settings.Conf.RedisConfig)

}
