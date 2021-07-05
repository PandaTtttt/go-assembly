package postutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PandaTtttt/go-assembly/util/m"
	"github.com/PandaTtttt/go-assembly/util/must"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	jsonType       = "application/json"
	urlencodedType = "application/x-www-form-urlencoded"
)

// PostClient provides a convenient way to send a post request.
// For some simple request which the caller only want the response body
// and don't care about the response header or status, ResponseBody will
// be a good choice to fetch the body and close TCP connection in a implicit way.

type PostClient struct {
	*http.Client
}

type headerField struct {
	key   string
	value string
}

func Client() *PostClient {
	return ClientInherit(nil)
}

func ClientInherit(cli *http.Client) *PostClient {
	// If nil, http.DefaultClient is used.
	if cli == nil {
		return &PostClient{
			Client: http.DefaultClient,
		}
	}
	return &PostClient{
		Client: cli,
	}
}

func Header(p m.M) []*headerField {
	var res []*headerField
	for k, v := range p {
		res = append(res, &headerField{
			key:   k,
			value: v.(string),
		})
	}
	return res
}

type response struct {
	res *http.Response
	err error
}

func (r *response) Response() (*http.Response, error) {
	return r.res, r.err
}

func (r *response) ResponseBody() ([]byte, error) {
	defer func() {
		if r.res == nil {
			return
		}
		must.Close(r.res.Body)
	}()

	if r.res == nil {
		return nil, r.err
	}
	body := must.Byte(ioutil.ReadAll(r.res.Body))
	if r.res.StatusCode >= http.StatusBadRequest {
		logger().Println(r.res.StatusCode, string(body))
	}
	return body, r.err
}

func (c *PostClient) RetryPostWithJsonType(params m.M,
	url string, retries int, header ...*headerField) *response {
	rawJson := must.Byte(json.Marshal(&params))
	header = append(header, &headerField{
		key:   "Content-Type",
		value: jsonType,
	})
	return c.postInternal(bytes.NewBuffer(rawJson), url, retries, header...)
}

func (c *PostClient) PostWithJsonType(params m.M,
	url string, header ...*headerField) *response {
	return c.RetryPostWithJsonType(params, url, 0, header...)
}

func (c *PostClient) RetryPostWithUrlencoded(params m.M,
	url string, retries int, header ...*headerField) *response {
	body := strings.NewReader(formatPostStr(params))
	header = append(header, &headerField{
		key:   "Content-Type",
		value: urlencodedType,
	})
	return c.postInternal(body, url, retries, header...)
}

func (c *PostClient) PostWithUrlencoded(params m.M,
	url string, header ...*headerField) *response {
	return c.RetryPostWithUrlencoded(params, url, 0, header...)
}

func (c *PostClient) RetryPost(url, contentType string, body io.Reader, retries int) *response {
	return c.postInternal(body, url, retries, Header(m.M{"Content-Type": contentType})...)
}

func (c *PostClient) Post(url, contentType string, body io.Reader) *response {
	return c.postInternal(body, url, 0, Header(m.M{"Content-Type": contentType})...)
}

func (c *PostClient) Do(req *http.Request) *response {
	res, err := c.Client.Do(req)
	if err != nil {
		logger().Println(err)
	}
	return &response{
		res: res,
		err: err,
	}
}

func (c *PostClient) postInternal(body io.Reader, url string, retries int, header ...*headerField) *response {
	// since http.Client.Do always close request body even on errors,
	// we create new request on every retry.
	for {
		request, err := http.NewRequest("POST", url, body)
		if err != nil {
			return &response{err: err}
		}
		for _, v := range header {
			// aims to skip https://github.com/golang/go/issues/7682
			if v.key == "Host" {
				request.Host = v.value
			}
			request.Header.Set(v.key, v.value)
		}
		res, err := c.Do(request).Response()
		if err != nil && retries != 0 {
			time.Sleep(time.Millisecond * 500)
			retries--
			continue
		}
		return &response{
			res: res,
			err: err,
		}
	}
}

// All errors and bad response have been recorded into (exec path)/posterr.log by default.
var _logger *log.Logger
var once sync.Once

func logger() *log.Logger {
	if _logger == nil {
		once.Do(func() {
			out, err := os.OpenFile("./posterr.log",
				os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
			must.Must(err)
			_logger = log.New(out, "", log.LstdFlags)
		})
	}
	return _logger
}

func SetErrOutput(file string) {
	must.Must(os.MkdirAll(filepath.Dir(file), os.ModePerm))
	out, err := os.OpenFile(file,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	must.Must(err)
	_logger = log.New(out, "", log.LstdFlags)
}

func formatPostStr(data m.M) string {
	postStr := ""
	for k, v := range data {
		postStr += fmt.Sprintf("%v=%v&", k, v)
	}
	return postStr
}
