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

	"go.uber.org/zap"

	"github.com/x-hezhang/gowebapp/routes"

	"github.com/x-hezhang/gowebapp/dao/redis"

	"github.com/x-hezhang/gowebapp/dao/mysql"

	"github.com/x-hezhang/gowebapp/logger"

	"github.com/x-hezhang/gowebapp/settings"
)

// Go Web较通用的脚手架模板

func main() {
	// 配置初始化
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

	defer func() { _ = zap.L().Sync() }()

	// MySQL初始化
	if err := mysql.Init(); err != nil {
		log.Fatalf("init database failed! %v\n", err)
	} else {
		fmt.Println("init database success!")
	}

	defer mysql.Close()

	// Redis初始化
	if err := redis.Init(); err != nil {
		log.Fatalf("init redis failed! %v\n", err)
	} else {
		fmt.Println("init redis success!")
	}

	defer redis.Close()

	// 注册路由
	r := routes.Init()
	routes.RegisterMiddleware(r, logger.GinLogger(), logger.GinRecovery(true))

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
