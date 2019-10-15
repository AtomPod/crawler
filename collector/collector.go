package collector

//Collector 数据收集器
type Collector interface {
	Collect(data interface{}) error
	Close() error
}
