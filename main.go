package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/x-hezhang/gowebapp/app/index"

	"go.uber.org/zap"

	"github.com/x-hezhang/gowebapp/routes"

	"github.com/x-hezhang/gowebapp/dao/redis"

	"github.com/x-hezhang/gowebapp/dao/mysql"

	"github.com/x-hezhang/gowebapp/logger"

	"github.com/x-hezhang/gowebapp/settings"
)

// Go Web较通用的脚手架模板

func main() {
	// 解析用户指定的配置文件
	configPath := settings.Parse()

	// 配置初始化
	if err := settings.Init(configPath); err != nil {
		log.Fatalf("init settings failed! %v\n", err)
	} else {
		log.Println("init settings success!")
	}

	// 日志初始化
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		log.Fatalf("init logger failed! %v\n", err)
	} else {
		log.Println("init logger success!")
	}

	defer func() { _ = zap.L().Sync() }()

	// MySQL初始化
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		log.Fatalf("init database failed! %v\n", err)
	} else {
		log.Println("init database success!")
	}

	defer mysql.Close()

	// Redis初始化
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		log.Fatalf("init redis failed! %v\n", err)
	} else {
		log.Println("init redis success!")
	}

	defer redis.Close()

	// 设置运行模式
	gin.SetMode(settings.Conf.AppConfig.Mode)

	// 路由初始化
	r := routes.Init()

	// 注册中间件
	routes.RegisterMiddlewares(
		r,
		logger.GinLogger(),
		logger.GinRecovery(true),
	)

	// 注册路由
	v1 := r.Group("/api/v1")
	routes.RegisterRoutes(v1, index.Routes)

	// 服务启动和优雅关闭
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.AppConfig.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server forced to shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
