# 简介
> * 长期总结并维护golang项目中的通用基础包。
> * 中的代码会被很多个不同项目引用，所以代码的修改和二次开发需要额外小心，通常情况下不应该直接修改原有代码。
> * 除`proto`外的任何修改需要单独的开发分支，开发完成后提交merge request，在审核无误之后方可merge。

### simplejson
* `simplejson`可有效处理多层嵌套结构的json数据(常见于非关系型数据库，ES等)，无需自定义struct。
* `simplejson.go` fork自[simplejson](https://github.com/bitly/go-simplejson), 参照openPR整体API有所调整，并修复了一些issue。
* `sql_datatype.go` 实现 [gromV2](https://gorm.io/docs/v2_release_note.html) 接口， `simplejson.JSON`类型直接映射到数据库模型（可在orm中使用）。
* json的`set`和`get`方法均为O(1)，同时提供数据库json操作API。
* 具体使用细节详见包内测试文件。

### zlog
* `zlog`封装了高性能日志库 [zap](https://github.com/uber-go/zap), 性能表现优异。
* 在 [zap](https://github.com/uber-go/zap) 的基础上提供了日志文件分割备份，定期删除。
* 提供 gin 框架支持，可用作 gin 的日志输出。
* 可配置多个输出流，具备日志消息直传kafka的能力。
* 具体使用细节详见`examples/zlog`。

### util
`util`中的方法应该是普遍适用的、无状态的。<br>
原则上`util`不应引用项目中的其他包，如果一个函数实现必须要引用项目其他包，那这个函数就不应该放在`util`下。
* `util/m` 是`map[string]interface{}`的简单表示，此数据类型用途广泛但写法繁琐，多数情况下使用`m.M{"field":"value"}`即可。
* `util/must` 为`panic`提供了更优雅的写法，使用时请遵循`gotips.md`中的`何时panic?`原则。

### queue
* `queue`是一个线程安全的任务队列，具体使用见包内测试文件。<br>

### postutil
* `postutil`为发起post请求提供了封装，具体使用细节详见`examples/postuitl`。<br>
* `postutil`默认会将请求过程中出现的错误打印在`(exec path)/posterror.log`中。<br>
* `postutil`可通过`SetErrOutput`重新指定log输出(在`main`包下`init`中使用）。<br>

### errs
`errs`提供自定错误类型和类型判断。<br>
`errs`包含八种通用的错误类型，在此之上可添加业务逻辑专属的错误类型。<br>
具体使用细节见`examples/errs`。
```go
const (
	// Internal is the generic error that maps to HTTP 500.
	Internal RetCode = iota + 100001
	// NotFound indicates a given resource is not found.
	NotFound
	// Forbidden indicates the user doesn't have the permission to
	// perform given operation.
	Forbidden
	// Unauthenticated indicates the oauth2 authentication failed.
	Unauthenticated
	// InvalidArgument indicates the input is invalid.
	InvalidArgument
	// InvalidConfig indicates the config is invalid.
	InvalidConfig
	// Conflict indicates a database transactional conflict happens.
	Conflict
	// TryAgain indicates a temporary outage and retry
	// could eventually lead to success.
	TryAgain
)
```

### env
`env`的作用是处理项目初始化时的必要依赖，通过环境变量或者手动传参的形式进行相关配置。
* `env/env.go`提供环境变量的解析。
* `env/mysql`提供mysql配置的初始化。
* `env/redis`提供redis配置的初始化。
* `env/server`提供gin server配置的初始化。
* `env/conn`提供mysql和redis的连接实体，生产环境下mysql配置慢查询日志。

目前可供配置的环境变量：
```go
const (
	DBDialect            = "DB_DIALECT"
	DBDatabase           = "DB_DATABASE"
	DBUser               = "DB_USER"
	DBPassword           = "DB_PASSWORD"
	DBCharset            = "DB_CHARSET"
	DBHost               = "DB_HOST"
	DBPort               = "DB_PORT"
	DBMaxIdleConns       = "DB_MAXIDLECONNS"
	DBMaxOpenNConns      = "DB_MAXOPENCONNS"

	RedisSentinel        = "REDIS_SENTINEL"
	RedisMaster          = "REDIS_MASTERNAME"
	RedisNetwork         = "REDIS_NETWORK"
	RedisAddr            = "REDIS_ADDR"
	RedisPassword        = "REDIS_PASSWORD"

	GoUseDefault         = "GO_ENVUSEDEFAULT"
	GoEnv                = "GO_ENV"
	GoLogDir             = "GO_LOGDIR"
	GoApiPrefix          = "GO_APIPREFIX"
	GoMaxMultipartMemory = "GO_MAXMUlTIPARTMEMORY"
	GoPort               = "GO_PORT"
)
```
具体使用见`examples/env`。

### cache
* `cache`是并发安全的，具体使用见包内测试文件。<br>
* `cache`灵感来源于《The Go Programming Language》p212-218。<br>
* `cache`主要解决[函数记忆](https://en.wikipedia.org/wiki/Memoization)问题，即缓存函数的结果，达到多次调用只需计算一次的效果。<br>
* `cache`拥有GC，可以用来做redis的二级缓存（无lease，会有缓存一致性问题，酌情使用）。

### atomic
`atomic` 封装了`sync/atomic`，只是提供了更简便的写法，并无任何功能上的改进。灵感来源：`github.com/uber-go/atomic`。<br>
需要注意的是原子操作属于`low-level`编程，大量的原子操作逻辑复杂且容易出错，并发编程请优先使用`channel`和`sync`包下的方法。

### api
`api`为`gin`框架提供一些辅助函数。目前包括请求头校验中间件和encode返回结果函数。





