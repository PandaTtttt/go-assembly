package conn

import (
	"github.com/PandaTtttt/go-assembly/env"
	"github.com/PandaTtttt/go-assembly/env/mysql"
	"github.com/PandaTtttt/go-assembly/util/must"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

// DB 数据库连接
var _DB *gorm.DB
var dbOnce sync.Once

// DB follows singleton pattern
func DB() *gorm.DB {
	dbOnce.Do(initDBConn)
	return _DB
}

func initDBConn() {
	var logger gormlogger.Interface
	goEnv := os.Getenv(env.GoEnv)
	if goEnv == env.DevelopMode || goEnv == "" {
		logger = gormlogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gormlogger.Config{
			SlowThreshold: 100 * time.Millisecond,
			LogLevel:      gormlogger.Info,
			Colorful:      true,
		})
	} else {
		logger = gormlogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gormlogger.Config{
			// silent in product env.
			LogLevel: gormlogger.Silent,
		})
	}

	db, err := gorm.Open(driver.Open(mysql.Config.URL), &gorm.Config{Logger: logger})
	must.Must(err)

	sqlDB, err := db.DB()
	must.Must(err)
	sqlDB.SetMaxIdleConns(mysql.Config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysql.Config.MaxOpenConns)

	_DB = db
}
