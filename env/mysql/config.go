package mysql

import (
	"fmt"
	"github.com/PandaTtttt/go-assembly/env"
	"strconv"
)

type config struct {
	Dialect      string
	Database     string
	User         string
	Password     string
	Host         string
	Port         int
	Charset      string
	URL          string
	MaxIdleConns int
	MaxOpenConns int
}

// Config 数据库相关配置
var Config config

type Params struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

// Init populates Config by environment variables or given params and default value.
func Init(p *Params) {
	var err error
	Config.Dialect = env.Get(env.DBDialect, "mysql")
	Config.Database = env.Get(env.DBDatabase, p.Database)
	Config.User = env.Get(env.DBUser, p.User)
	Config.Password = env.Get(env.DBPassword, p.Password)
	Config.Host = env.Get(env.DBHost, p.Host)
	Config.Charset = env.Get(env.DBCharset, "utf8mb4")

	Config.Port, err = strconv.Atoi(env.Get(env.DBPort, p.Port))
	if err != nil {
		panic(env.Invalid(env.DBPort))
	}

	Config.MaxIdleConns, err = strconv.Atoi(env.Get(env.DBMaxIdleConns, "300"))
	if err != nil {
		panic(env.Invalid(env.DBMaxIdleConns))
	}
	Config.MaxOpenConns, err = strconv.Atoi(env.Get(env.DBMaxOpenNConns, "900"))
	if err != nil {
		panic(env.Invalid(env.DBMaxOpenNConns))
	}

	Config.URL = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		Config.User, Config.Password, Config.Host, Config.Port, Config.Database, Config.Charset)
}
