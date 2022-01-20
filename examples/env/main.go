package main

import (
	"github.com/PandaTtttt/go-assembly/env/mysql"
	"github.com/PandaTtttt/go-assembly/env/redis"
	"github.com/PandaTtttt/go-assembly/env/server"
)

func main() {
	// 初始化mysql，配置参数优先走环境变量，其次才是mysql.Params
	mysql.Init(&mysql.Params{
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "xxx",
		User:     "root",
		Password: "",
	})
	// 初始化后方可使用全局的DB连接实体conn.DB()
	// 在调用mysql.Init之前调用conn.DB()会panic
	// example: conn.DB().New()

	// 初始化redis，配置参数优先走环境变量，其次才是redis.Params
	redis.Init(&redis.Params{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       "0",
	})
	// 初始化后方可使用全局的RS连接实体conn.RS()
	// 在调用redis.Init之前调用conn.RS()会panic
	// example: conn.RS().Do()

	// 初始化gin server，配置参数优先走环境变量，其次才是server.Params
	server.Init(&server.Params{
		ApiPrefix: "v1",
		Port:      "8023",
	})

	// 以上的初始化只应该写在main函数里也可以写在项目中的其他init方法中。
}
