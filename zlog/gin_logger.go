package zlog

import (
	"fmt"
	"github.com/PandaTtttt/go-assembly/util/must"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

func GinLogger(conf Config) (info gin.HandlerFunc, recovery gin.HandlerFunc) {
	encoderConf := zapcore.EncoderConfig{
		TimeKey:        zapcore.OmitKey,
		LevelKey:       zapcore.OmitKey,
		NameKey:        zapcore.OmitKey,
		CallerKey:      zapcore.OmitKey,
		MessageKey:     zapcore.OmitKey,
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			nano := t.UnixNano()
			milli := nano / int64(time.Millisecond)
			enc.AppendInt64(milli)
		},
	}

	logger, err := New(conf, encoderConf, zap.AddStacktrace(zap.WarnLevel))
	must.Must(err)
	return func(c *gin.Context) {

			start := time.Now()
			timeLocal := start.Format("02/Jan/2006:15:04:05 +0800")

			c.Next()

			l := float64(time.Since(start)) / float64(time.Millisecond)
			latency := fmt.Sprintf("%.3f", l)

			logger.Info(c.ClientIP(),
				zap.String("log_id", c.GetString("request_id")),
				zap.String("time_local", timeLocal),
				zap.String("remote_addr", c.ClientIP()),
				zap.String("http_referer", c.Request.Referer()),
				zap.String("request", fmt.Sprintf("%s %s %s", c.Request.Method, c.Request.RequestURI, c.Request.Proto)),
				zap.Int("status", c.Writer.Status()),
				zap.Int("body_bytes_sent", c.Writer.Size()),
				zap.String("http_user_agent", c.Request.UserAgent()),
				zap.String("http_x_forwarded_for", c.ClientIP()),
				zap.String("request_method", c.Request.Method),
				zap.String("uri", c.Request.RequestURI),
				zap.String("request_time", latency),
			)
		},
		func(c *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					// Check for a broken connection, as it is not really a
					// condition that warrants a panic stack trace.
					var brokenPipe bool
					if ne, ok := err.(*net.OpError); ok {
						if se, ok := ne.Err.(*os.SyscallError); ok {
							if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
								strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
								brokenPipe = true
							}
						}
					}

					now := time.Now().UnixNano() / int64(time.Millisecond)
					httpRequest, _ := httputil.DumpRequest(c.Request, false)

					if brokenPipe {
						logger.Error(c.Request.URL.Path,
							zap.Int64("timestamp", now),
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
						// If the connection is dead, we can't write a status to it.
						c.Error(err.(error)) // nolint: errcheck
						c.Abort()
						return
					}

					logger.Error("Recovery from panic",
						zap.Int64("timestamp", now),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}()
			c.Next()
		}
}
