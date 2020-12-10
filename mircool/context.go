package mircool

import (
	json "encoding/json"
	"fmt"
	"io"
	"net/http"
)

type M map[string]interface{}

type Context struct {
	Resp       http.ResponseWriter //响应
	Req        *http.Request       //请求
	Path       string              //请求路径
	Method     string              //请求方式
	Params     map[string]string   //请求参数（url或者formData参数）
	StatusCode int                 //响应状态码
	handlers   []HandleFunc        //中间件
	index      int
}

func newContext(resp http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Resp:   resp,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	i := len(c.handlers)
	for ; c.index < i; c.index++ {
		c.handlers[c.index](c)
	}

}

func (c *Context) Param(key string) string {
	s := c.Params[key]
	return s
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Resp.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.Resp.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Resp.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Resp)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Resp, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) BindJson(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&obj); err != nil {
		return err
	}
	return nil
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Resp.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Resp.Write([]byte(html))
}
