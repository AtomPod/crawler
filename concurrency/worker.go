package concurrency

import (
	"context"
	"log"
)

//Worker 工作者
type Worker struct {
	d      *Dispatcher
	ctx    context.Context
	cancel context.CancelFunc
	jobs   chan Job
}

//NewWorker 创建一个Worker
func NewWorker(pctx context.Context, d *Dispatcher) *Worker {
	if pctx == nil {
		pctx = context.Background()
	}
	ctx, cancel := context.WithCancel(pctx)

	return &Worker{
		d:      d,
		ctx:    ctx,
		cancel: cancel,
		jobs:   make(chan Job),
	}
}

//Run 执行worker
func (w *Worker) Run() {
	go func() {
		w.work()
	}()
}

//Process 执行Job
func (w *Worker) Process(j Job) {
	select {
	case w.jobs <- j:
	case <-w.ctx.Done():
		return
	}
}

func (w *Worker) work() {
	for {
		w.d.workerPool <- w

		select {
		case job := <-w.jobs:
			if err := job.Execute(); err != nil {
				log.Println(err)
			}
		case <-w.ctx.Done():
			return
		}
	}
}

//Close 关闭工作者
func (w *Worker) Close() error {
	w.cancel()
	return nil
}
