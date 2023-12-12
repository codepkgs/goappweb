package logger

import (
	"fmt"
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

func getConsoleEncoder() zapcore.Encoder {
	encoder := zap.NewProductionEncoderConfig()
	encoder.TimeKey = "time"
	encoder.MessageKey = "message"
	encoder.EncodeTime = zapcore.RFC3339TimeEncoder
	encoder.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoder.EncodeCaller = zapcore.ShortCallerEncoder
	encoder.EncodeDuration = zapcore.SecondsDurationEncoder

	return zapcore.NewConsoleEncoder(encoder)
}

func getFileWriterSyncer(cfg *settings.LogConfig) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
	})
}

var ErrUnsupportedOutput = fmt.Errorf("unsupported logger output, supported outputs: file, stdout")

func Init(cfg *settings.LogConfig) (err error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return
	}

	var core zapcore.Core
	switch cfg.Output {
	case "file":
		core = zapcore.NewCore(getJsonEncoder(), getFileWriterSyncer(cfg), level)
	case "stdout":
		core = zapcore.NewCore(getConsoleEncoder(), zapcore.AddSync(os.Stdout), level)
	default:
		return ErrUnsupportedOutput
	}

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
