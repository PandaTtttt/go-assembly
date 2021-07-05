# Go编程指南

官方指导(建议精读)：
* [Effective Go](https://golang.org/doc/effective_go.html)
* [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

民间指导：
* [go101](https://github.com/go101/go101)

以下是一些重要原则的重申及额外原则，以及一些编码小tips，任何人都可以在此分享自己的心得。

### 何时panic？
任何goroutine中的`panic`都会导致整个application退出，此处也包括`util/must`中的所有方法，在使用`must.xxx()`时一定要遵循以下原则：
* panic可以发生在程序初始化时间。
* 初始化完成以后，panic不能让服务器崩溃。
* panic只能是由于两种事情发生：
    1. 错误的原因是人为的判断错误，比如程序员相信某个`struct`可以`marshal`成json，不可能出错，那么万一出错，可以panic。
    1. 错误是由于核心服务出错，比如从mysql查询结果，服务器返回异常结果，这个时候可以panic。
* panic不能由于以下原因：
    1. 数据录入不正确。
    1. 最终用户提供的参数不正确。
    1. 任何其他不是程序员直接引发的错误。
* 总体而言，导致panic的错误必须是程序员可以修复的错误。
* 如果你不确定，不要panic。

### 不要吃掉错误
如果一个函数返回有错误，你必须选择其中一种：
1. 把错误作为结果从函数返回。
1. 把错误输出到日志文件中。
1. panic，见上文讨论何时panic。
1. 如果错误是预计得到的，正确检查以后，进行相应处理，比如重试。

以下处理方法就是吃掉错误：
1. 完全不检查错误。
1. `_ = getError()`，直接忽视错误。
1. `fmt.Print(err)`，然后不管。

### err处理流程
不应该出现的流程：
```go
if err != nil {
	// error handling
} else {
	// normal code
}
```
正确的流程：
```go
if err != nil {
	// error handling
	return // or continue, etc.
}
// normal code
```

### 不要返回 in-band Error
in-band Error指函数没有按照预期执行或出现错误时，通过返回`-1`或`""`等值来表示错误。<br>
由于golang支持多个返回值，请用额外的返回值来指示函数是否成功执行。<br>

不应该出现的函数签名：

```go
// Lookup returns the value for key or "" if there is no mapping for key.
func Lookup(key string) string
```
正确的函数签名以及流程如下：
```go
// Lookup returns the value for key or ok=false if there is no mapping for key.
func Lookup(key string) (value string, ok bool)

value, ok := Lookup(key)
if !ok {
	return fmt.Errorf("no value for %q", key)
}
return Parse(value)
```

### 空slice声明
当你要创建一个slice，并不知道它的长度的话使用如下风格：
```go
var someSlice []MyStruct
for ... {
  someSlice = append(someSlice, ...)
}
```
不要使用```someSlice := []MyStruct{}```或者```someSlice := make(...)```（除非你能预判大小）。

这样的好处是如果这个slice最后是0长度（没有任何append），我们省了一次内存分配。

### 函数参数太多的时候，考虑写Params struct。
例如
```go
func SetUser(name, gender, birthday, address, phone, email string, age int) { ... }
```
这个函数接受很多参数，而且类型都一样。这样调用的时候，参数的位置很容易搞错，也不容易读懂。
可以改成
```go
type SetUserParam struct {
	name     string
	gender   string
	birthday string
	address  string
	phone    string
	email    string
	age      int
}
func SetUser(p SetUserParam) { ... }
```

### 关于资源泄漏

对于大部分实现了`io.Closer`接口的对象，都标志着其生命周期中占据了一定的系统资源，使用完毕后需要手动调用`Close`方法释放资源，否则会造成资源泄漏。

常见如下：
* `os.File` 占据一个文件句柄。
* `http.Response.Body` 占据一个socket连接。
* `sql.Stmt` 占据一个prepared语句。

### channel关闭原则
`channel`的关闭与上述的资源泄漏没有关系，一个`channel`的关闭本质上是一种广播行为，程序员在认为需要关闭的时候才进行关闭，并遵循以下原则：
* 不要在数据接收方进行关闭。
* 有多个并行的数据发送者时不要进行关闭。
* 只有`channel`唯一的发送者才能进行关闭。

`channel`关闭原则本质上是为了防止重复关闭一个已经关闭的`channel`。

### 竞态检测
在开发阶段，我们应该总是使用`go run/build -race`来帮助检测是否有竞态产生。<br>
`-race`不能百分百保证程序中没有竞态，检测依赖于程序的具体执行。<br>
不要忽视`-race`带来的警告，即使你认为你的程序在并发读写上没有问题，但有竞态的代码可能会导致编译器的一些错误行为。<br>
延伸阅读 [Does the Go race detector catch all data race bugs?](https://medium.com/@val_deleplace/does-the-race-detector-catch-all-data-races-1afed51d57fb)




