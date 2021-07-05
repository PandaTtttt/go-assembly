package env

import (
	"fmt"
	"os"
)

const (
	DevelopMode = "develop"
	TestMode    = "test"
	ProductMode = "product"
)

//声明环境变量常量
const (
	DBDialect       = "DB_DIALECT"
	DBDatabase      = "DB_DATABASE"
	DBUser          = "DB_USER"
	DBPassword      = "DB_PASSWORD"
	DBCharset       = "DB_CHARSET"
	DBHost          = "DB_HOST"
	DBPort          = "DB_PORT"
	DBMaxIdleConns  = "DB_MAXIDLECONNS"
	DBMaxOpenNConns = "DB_MAXOPENCONNS"

	RedisSentinel = "REDIS_SENTINEL"
	RedisMaster   = "REDIS_MASTERNAME"
	RedisNetwork  = "REDIS_NETWORK"
	RedisAddr     = "REDIS_ADDR"
	RedisPassword = "REDIS_PASSWORD"

	GoUseDefault         = "GO_ENVUSEDEFAULT"
	GoEnv                = "GO_ENV"
	GoLogDir             = "GO_LOGDIR"
	GoApiPrefix          = "GO_APIPREFIX"
	GoMaxMultipartMemory = "GO_MAXMUlTIPARTMEMORY"
	GoPort               = "GO_PORT"
)

func Invalid(key string) string {
	return fmt.Sprintf("invalid: %s=%v", key, os.Getenv(key))
}

func Get(key, d string) string {
	res := os.Getenv(key)
	if res == "" && os.Getenv(GoUseDefault) != "false" {
		return d
	}
	return res
}
