package curl

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Request构造类
type Request struct {
	cli      *http.Client
	req      *http.Request
	Raw      *http.Request
	Method   string
	Url      string
	Headers  map[string]string
	Cookies  map[string]string
	Queries  map[string]string
	PostData map[string]interface{}
}

// 创建一个Request实例
func NewRequest() *Request {
	return &Request{}
}

// 设置请求方法
func (this *Request) SetMethod(method string) *Request {
	this.Method = method
	return this
}

// 设置请求地址
func (this *Request) SetUrl(url string) *Request {
	this.Url = url
	return this
}

// 设置请求头
func (this *Request) SetHeaders(headers map[string]string) *Request {
	this.Headers = headers
	return this
}

// 将用户自定义请求头添加到http.Request实例上
func (this *Request) setHeaders() error {
	for k, v := range this.Headers {
		this.req.Header.Set(k, v)
	}
	return nil
}

// 设置请求cookies
func (this *Request) SetCookies(cookies map[string]string) *Request {
	this.Cookies = cookies
	return this
}

// 将用户自定义cookies添加到http.Request实例上
func (this *Request) setCookies() error {
	for k, v := range this.Cookies {
		this.req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	return nil
}

// 设置url查询参数
func (this *Request) SetQueries(queries map[string]string) *Request {
	this.Queries = queries
	return this
}

// 将用户自定义url查询参数添加到http.Request上
func (this *Request) setQueries() error {
	q := this.req.URL.Query()
	for k, v := range this.Queries {
		q.Add(k, v)
	}
	this.req.URL.RawQuery = q.Encode()
	return nil
}

// 设置post请求的提交数据
func (this *Request) SetPostData(postData map[string]interface{}) *Request {
	this.PostData = postData
	return this
}

// 发起get请求
func (this *Request) Get() (*Response, error) {
	return this.Send(this.Url, http.MethodGet)
}

// 发起post请求
func (this *Request) Post() (*Response, error) {
	return this.Send(this.Url, http.MethodPost)
}

// 发起请求
func (this *Request) Send(url string, method string) (*Response, error) {
	// 初始化Response对象
	response := NewResponse()

	// 初始化http.Client对象
	this.cli = &http.Client{Timeout: 30 * time.Second}

	// 检测请求url是否填了
	if url == "" {
		panic("Lack of request url")
	}

	// 检测请求方式是否填了
	if method == "" {
		panic("Lack of request method")
	}

	// 加载用户自定义的post数据到http.Request
	var payload io.Reader
	if method == "POST" && this.PostData != nil {
		if jData, err := json.Marshal(this.PostData); err != nil {
			panic(err)
		} else {
			payload = bytes.NewReader(jData)
		}
	} else {
		payload = nil
	}

	if req, err := http.NewRequest(method, url, payload); err != nil {
		panic(err)
	} else {
		this.req = req
	}

	this.setHeaders()
	this.setCookies()
	this.setQueries()

	this.Raw = this.req

	if resp, err := this.cli.Do(this.req); err != nil {
		panic(err)
	} else {
		response.Raw = resp
	}

	defer response.Raw.Body.Close()

	response.parseHeaders()
	response.parseBody()

	return response, nil
}
