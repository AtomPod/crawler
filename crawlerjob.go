package crawler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/phantom-atom/crawler/parser"
)

var (
	//bytes.Buffer对象池
	bufPool = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

type crawlerJob struct {
	e *Engine
	r *parser.Request
}

func (c *crawlerJob) Execute() error {
	if c.r.WaitGroup != nil {
		defer c.r.WaitGroup.Done()
	}

	buf, err := c.get()

	if err != nil {
		return err
	}

	result, err := c.r.Parser.Parse(c.r, buf)
	bufPool.Put(buf)

	if err != nil {
		return err
	}

	if result != nil {
		requests := result.Requests()
		for _, req := range requests {
			c.e.Process(req)
		}
	}
	return nil
}

func (c *crawlerJob) get() (*bytes.Buffer, error) {
	resp, err := http.Get(c.r.URL.String())

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	//如果请求的状态码不是OK，那么需要将resp.Body
	//清空，以便复用连接，并返回错误信息
	if resp.StatusCode != http.StatusOK {
		io.Copy(ioutil.Discard, resp.Body)
		return nil, fmt.Errorf("request: url=%s, status_code=%d", c.r.URL.String(), resp.StatusCode)
	}

	//从对象池中获取一个Buffer
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	_, err = io.Copy(buf, resp.Body)

	if err != nil {
		bufPool.Put(buf)
		return nil, err
	}
	return buf, nil
}
