package parser

import (
	"net/url"
	"sync"

	"github.com/phantom-atom/crawler/collector"
)

//Request 请求结构体
type Request struct {
	//请求的URL
	URL *url.URL

	//解析器，用于解析URL请求后的数据
	Parser Parser

	//等待组，用于等待解析完成，该变量
	//可以通过每次解析的时候，传递给下一个Request
	//这样只要在最开始调用的时候传入，可以等待所有请求结束
	WaitGroup *sync.WaitGroup

	//数据收集器，可用于数据接收使用
	Collector collector.Collector
}
