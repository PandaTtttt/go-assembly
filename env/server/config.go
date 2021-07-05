package server

import (
	"github.com/PandaTtttt/go-assembly/env"
	"github.com/PandaTtttt/go-assembly/util/must"
	"os"
	"path/filepath"
	"strconv"
)

type config struct {
	LogDir             string
	LogFile            string
	APIPrefix          string
	MaxMultipartMemory int
	Port               int
}

// Config 服务器相关配置
var Config config

type Params struct {
	ApiPrefix string
	Port      string
}

// Init populates Config by environment variables or given params and default value.
func Init(p *Params) {
	var err error
	Config.LogDir = env.Get(env.GoLogDir, must.String(os.Getwd()))
	Config.APIPrefix = env.Get(env.GoApiPrefix, p.ApiPrefix)
	Config.MaxMultipartMemory, err = strconv.Atoi(env.Get(env.GoMaxMultipartMemory, "3"))
	if err != nil {
		panic(env.Invalid(env.GoMaxMultipartMemory))
	}
	Config.Port, err = strconv.Atoi(env.Get(env.GoPort, p.Port))
	if err != nil {
		panic(env.Invalid(env.GoPort))
	}

	Config.LogFile = filepath.Join(Config.LogDir, "gin.log")
}
