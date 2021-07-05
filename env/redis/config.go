package redis

import (
	"github.com/PandaTtttt/go-assembly/env"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

type config struct {
	Sentinel   bool
	MasterName string
	Network    string
	Addr       string
	Password   string
}

// Config redis相关配置
var Config config
var Option interface{}

type Params struct {
	Addr     string
	Password string
}

// Init populates Config by environment variables or given params and default value.
func Init(p *Params) {
	var err error
	Config.Sentinel, err = strconv.ParseBool(env.Get(env.RedisSentinel, "false"))
	if err != nil {
		panic(env.Invalid(env.RedisSentinel))
	}
	Config.MasterName = env.Get(env.RedisMaster, "")
	Config.Network = env.Get(env.RedisNetwork, "tcp")
	Config.Addr = env.Get(env.RedisAddr, p.Addr)
	Config.Password = env.Get(env.RedisPassword, p.Password)
	if Config.Sentinel {
		Option = redis.FailoverOptions{
			MasterName:         Config.MasterName,
			SentinelAddrs:      []string{Config.Addr},
			Password:           Config.Password,
			MinIdleConns:       1,
			IdleTimeout:        time.Minute * 10,
			IdleCheckFrequency: 7 * time.Minute,
		}
	} else {
		Option = redis.Options{
			Network:            Config.Network,
			Addr:               Config.Addr,
			Password:           Config.Password,
			MinIdleConns:       1,
			IdleTimeout:        time.Minute * 10,
			IdleCheckFrequency: 7 * time.Minute,
		}
	}
}
