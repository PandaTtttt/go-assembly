package conn

import (
	rs "github.com/PandaTtttt/go-assembly/env/redis"
	"github.com/PandaTtttt/go-assembly/util/must"
	"github.com/go-redis/redis/v7"
	"sync"
)

var _RS *redis.Client
var rsOnce sync.Once

// RS follows singleton pattern
func RS() *redis.Client {
	rsOnce.Do(func() {
		initRedisConn(0)
	})
	return _RS
}

func initRedisConn(index int) {
	client := &redis.Client{}
	switch option := rs.Option.(type) {
	case redis.Options:
		option.DB = index
		client = redis.NewClient(&option)
	case redis.FailoverOptions:
		option.DB = index
		client = redis.NewFailoverClient(&option)
	}
	must.String(client.Ping().Result())
	_RS = client
}
