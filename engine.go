package crawler

import (
	"context"

	"github.com/phantom-atom/crawler/concurrency"
	"github.com/phantom-atom/crawler/parser"
)

//Engine 爬虫引擎
type Engine struct {
	dispatcher *concurrency.Dispatcher //调度器
}

//NewEngine 创建Engine，maxWorkers为同一时间段最大并发数，maxJobQueue为最大Job队列
func NewEngine(ctx context.Context, maxWorkers, maxJobQueue int) *Engine {
	return &Engine{
		dispatcher: concurrency.NewDispatcher(ctx, maxWorkers, maxJobQueue),
	}
}

//Run 运行解析
func (e *Engine) Run() {
	e.dispatcher.Run()
}

//Process 添加处理的请求
func (e *Engine) Process(req *parser.Request) {
	if req == nil {
		return
	}

	//如果等待组不为nil，那么添加一个等待信息
	//保证在req.WaitGroup.Wait等待之前已经添加一个等待信息
	if req.WaitGroup != nil {
		req.WaitGroup.Add(1)
	}

	//添加请求到调度器中
	e.dispatcher.Process(&crawlerJob{
		r: req,
		e: e,
	})
}

//Close 关闭
func (e *Engine) Close() error {
	return e.dispatcher.Close()
}
