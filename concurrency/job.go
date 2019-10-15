package concurrency

//Job 执行事情
type Job interface {
	Execute() error
}
