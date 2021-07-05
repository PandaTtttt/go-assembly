package main

import (
	"fmt"
	"github.com/PandaTtttt/go-assembly/api"
	"github.com/PandaTtttt/go-assembly/env"
	"github.com/PandaTtttt/go-assembly/env/server"
	"github.com/PandaTtttt/go-assembly/errs"
	"github.com/PandaTtttt/go-assembly/util/must"
	"github.com/PandaTtttt/go-assembly/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

// zlog的全局logger,即zlog.Info()...所使用的logger应该在main包下的init函数里进行初始化。
func init() {
	// Init所配置的输出源只在测试和生产环境起效，通过环境变量env.GoEnv设置，
	// 开发环境下所有日志输出到终端（stderr）
	zlog.Init(zlog.Config{
		// Info代表Info级别的日志输出源，可配置多个，同时输出。
		Info: []zlog.LogWriter{
			// FileWriter为文件输出源配置。
			&zlog.FileWriter{
				// 日志文件全路径
				File: "/zlog/info.log",
				// 备份目录
				BackupDir: "",
				// 文件备份/清空的阈值，单位:byte
				// 设置为零则不进行备份/清空操作
				MaxSize: 0,
				// 备份文件保留时间，单位:天
				MaxLifetime: 0,
				// 最大备份文件数量
				MaxBackups: 0,
				// 备份文件名是否使用UTC时间，默认为local
				UTCTime: false,
				// 是否启用清空功能代替备份功能
				Truncation: false,
			},
			// KafkaWriter为kafka输出源配置，日志消息会直传kafka。
			&zlog.KafkaWriter{
				Address: []string{"10.70.11.64:9092"},
				Topic:   "sptest",
			},
		},
		// Err代表 Warn及以上级别的日志输出源，可配置多个，同时输出。
		Err: []zlog.LogWriter{
			&zlog.FileWriter{
				File:       "/zlog/error.log",
				MaxSize:    100 << 20,
				Truncation: true,
			},
			&zlog.KafkaWriter{
				Address: []string{"10.70.11.64:9092"},
				Topic:   "sptest",
			},
		},
	})
}

// 此代码只做demo展示，不要尝试运行此文件。
func main() {
	// zlog 一共支持7种级别的日志输出，从Debug到Fatal为从小到大的级别排列。

	// Debug日志在生产和测试环境不会输出。
	zlog.Debug("debug msg", zap.String("id", "111"))

	// Info日志输出源通过zlog.Config.Info进行设置。
	zlog.Info("info msg", zap.String("id", "111"))

	// Warn及以下通过zlog.Config.Err进行设置。
	zlog.Warn("warn msg", zap.String("id", "111"))
	zlog.Error("error msg", zap.String("id", "111"))

	// DPanic 为development panic,调用此函数只有在开发环境下会引起程序panic。
	zlog.DPanic("dpanic msg", zap.String("id", "111"))

	zlog.Panic("panic msg", zap.String("id", "111"))
	zlog.Fatal("fatal msg", zap.String("id", "111"))

	// zlog.New可创建自定义的logger，如果不想使用全局logger的话。
	logger, err := zlog.New(zlog.Config{}, zapcore.EncoderConfig{})
	must.Must(err)
	logger.Info("xxx")

	// 任何实现io.Writer的实体都可以通过以下步骤设置成日志输出源，如果不想使用zlog.FileWriter或zlog.KafkaWriter的话。
	fd, err := os.Create("xxx")
	w := zlog.AddCallback(fd)
	_ = zlog.Config{
		Info: []zlog.LogWriter{w},
	}
}

func ginServer() {
	goEnv := os.Getenv(env.GoEnv)
	if goEnv != "" && goEnv != env.DevelopMode {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.New()
	app.Use(
		zlog.GinLogger(zlog.Config{
			// 定义gin的消息日志输出
			Info: []zlog.LogWriter{
				&zlog.FileWriter{
					// NOTE: server.Config在调用server.Init后才会进行初始化。
					File: server.Config.LogFile,
				},
			},
			// 定义gin的错误日志输出
			Err: []zlog.LogWriter{
				&zlog.FileWriter{
					File: filepath.Join(filepath.Dir(server.Config.LogFile), "errors", "error.log"),
				},
			},
		})...)
	route(app)

	must.Must(app.Run(":" + fmt.Sprintf("%d", server.Config.Port)))
}

func route(router *gin.Engine) {
	r := router.Group(server.Config.APIPrefix)
	{
		r.POST("/info", infoHandler)
		r.POST("/err", errHandler)
		r.POST("/panic", panicHandler)
	}
}

func infoHandler(c *gin.Context) {
	zlog.Info("this is a info msg", zap.String("id", "111"))
	api.ResultOK(c, nil)
}
func errHandler(c *gin.Context) {
	err := errs.InvalidArgument.New("this is a invalid argument error")
	zlog.Error(err.Error(), zap.String("id", "222"))
	api.Result400(c, err)
}
func panicHandler(c *gin.Context) {
	zlog.Panic("i'm panic!!!")
}
