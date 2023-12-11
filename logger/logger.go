package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/x-hezhang/gowebapp/settings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getJsonEncoder() zapcore.Encoder {
	encoder := zap.NewProductionEncoderConfig()
	encoder.TimeKey = "time"
	encoder.MessageKey = "message"
	encoder.EncodeTime = zapcore.RFC3339TimeEncoder
	encoder.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoder.EncodeCaller = zapcore.ShortCallerEncoder
	encoder.EncodeDuration = zapcore.SecondsDurationEncoder

	return zapcore.NewJSONEncoder(encoder)
}

func getFileWriterSyncer() zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   settings.Conf.LogConfig.Filename,
		MaxSize:    settings.Conf.LogConfig.MaxSize,
		MaxBackups: settings.Conf.LogConfig.MaxBackups,
		MaxAge:     settings.Conf.LogConfig.MaxAge,
	})
}

func Init() (err error) {
	level, err := zapcore.ParseLevel(settings.Conf.LogConfig.Level)
	if err != nil {
		return
	}
	core := zapcore.NewCore(getJsonEncoder(), getFileWriterSyncer(), level)
	lg := zap.New(core)
	zap.ReplaceGlobals(lg)
	return nil
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.String("ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("url", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.String("proto", c.Request.Proto),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("latency", cost),
		)
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
