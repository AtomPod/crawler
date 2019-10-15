package typed

import (
	"reflect"

	"github.com/phantom-atom/crawler/collector/typed/sender"
)

//Collector 类型实现
type Collector struct {
	s sender.Sender
}

//NewCollector 创建一个收集器
//目前medium仅支持chan类型
func NewCollector(medium interface{}) *Collector {
	c := &Collector{}
	c.initialize(medium)
	return c
}

func (c *Collector) initialize(v interface{}) {
	val := reflect.ValueOf(v)
	typ := val.Type()

	switch typ.Kind() {
	case reflect.Chan:
		c.s = sender.NewChannel(val, typ)
	default:
		panic("typed-collector: type is not supported")
	}
}

//Collect Collect接口实现
func (c *Collector) Collect(data interface{}) error {
	return c.s.Send(data)
}

//Close Close接口实现
func (c *Collector) Close() error {
	return c.s.Close()
}
