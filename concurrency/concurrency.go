package concurrency

import (
	"context"
)

//Dispatcher 调度器
type Dispatcher struct {
	workerPool chan *Worker
	jobQueue   chan Job
	ctx        context.Context
	cancel     context.CancelFunc
}

//NewDispatcher  创建调度器
func NewDispatcher(pctx context.Context, maxWorkers int, maxJobQueue int) *Dispatcher {
	if pctx == nil {
		pctx = context.Background()
	}
	ctx, cancel := context.WithCancel(pctx)

	return &Dispatcher{
		workerPool: make(chan *Worker, maxWorkers),
		jobQueue:   make(chan Job, maxJobQueue),
		ctx:        ctx,
		cancel:     cancel,
	}
}

//Run 执行调度器
func (d *Dispatcher) Run() {
	for i := 0; i < cap(d.workerPool); i++ {
		worker := NewWorker(d.ctx, d)
		worker.Run()
	}
	go d.dispatch()
}

//Process 执行JOB
func (d *Dispatcher) Process(j Job) {
	select {
	case <-d.ctx.Done():
	case d.jobQueue <- j:
	}
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			go func(j Job) {
				select {
				case w := <-d.workerPool:
					w.Process(j)
				case <-d.ctx.Done():
					return
				}
			}(job)
		case <-d.ctx.Done():
			return
		}
	}
}

//Close 关闭调度器
func (d *Dispatcher) Close() error {
	d.cancel()
	return nil
}
