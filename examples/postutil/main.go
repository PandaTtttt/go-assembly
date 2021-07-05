package main

import (
	"fmt"
	"github.com/PandaTtttt/go-assembly/postutil"
	"github.com/PandaTtttt/go-assembly/util/m"
	"github.com/PandaTtttt/go-assembly/util/must"
	"io/ioutil"
	"net/http"
	"time"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	must.Must(r.ParseForm())
	fmt.Println(r.PostForm)

	body := string(must.Byte(ioutil.ReadAll(r.Body)))
	if body != "" {
		fmt.Println(body)
	}

	fmt.Println(r.Header.Get("Content-Type"))
	fmt.Println(r.Header.Get("baz"))

	fmt.Fprintf(w, "hello world")
}

func main() {
	http.HandleFunc("/", Hello)
	go func() {
		must.Must(http.ListenAndServe("0.0.0.0:8800", nil))
	}()
	// wait for server
	time.Sleep(time.Second)

	// postutil会将请求过程中的错误记录在日志中，默认输出为os.stderr
	// SetErrOutput可以重定向输出到文件。
	postutil.SetErrOutput("/xxx/errs.log")

	// postutil.Client()默认继承http.DefaultClient
	c := postutil.Client()

	// c.Post 与http.DefaultClient.Post的不同之处在于后续提供两种链式调用，并且记录错误。
	_, _ = c.Post("http://127.0.0.1:8800", "application/json", nil).Response()

	// 拥有重试次数的post。
	_, _ = c.RetryPost("http://127.0.0.1:8800", "application/json", nil, 10).Response()

	// 不需要额外添加请求头的使用方法。
	_, _ = c.PostWithUrlencoded(m.M{"foo": "bar"}, "http://127.0.0.1:8800").Response()
	// POST
	// map[foo:[bar]]
	// application/x-www-form-urlencoded

	// 需要额外添加请求头的使用方法。
	c.PostWithUrlencoded(m.M{"foo": "bar"}, "http://127.0.0.1:8800", postutil.Header(m.M{"baz": "qux"})...)
	// console output as follow:
	// POST
	// map[foo:[bar]]
	// application/x-www-form-urlencoded
	// qux
	c.PostWithJsonType(m.M{"foo": "bar"}, "http://127.0.0.1:8800", postutil.Header(m.M{"baz": "qux"})...)
	// console output as follow:
	// POST
	// map[]
	// {"foo":"bar"}
	// application/json
	// qux

	// postutil.Client 可发起带有重试功能的请求。
	c.RetryPostWithJsonType(m.M{"foo": "bar"}, "http://127.0.0.1:8800", 10, postutil.Header(m.M{"baz": "qux"})...)
	c.RetryPostWithUrlencoded(m.M{"foo": "bar"}, "http://127.0.0.1:8800", 10, postutil.Header(m.M{"baz": "qux"})...)

	// 返回结果的两种获取方式
	// 1 标准返回格式，res中可获取所有的返回信息，调用者必须手动关闭tcp连接.
	res, err := c.PostWithUrlencoded(m.M{"foo": "bar"}, "http://127.0.0.1:8800").Response()
	must.Must(err)
	defer must.Close(res.Body)
	fmt.Println(res)

	// 2 只获取返回体，underlying的tcp连接已经在函数返回前关闭，调用者无需关心。
	body, err := c.PostWithUrlencoded(m.M{"foo": "bar"}, "http://127.0.0.1:8800").ResponseBody()
	fmt.Println(string(body))
}
