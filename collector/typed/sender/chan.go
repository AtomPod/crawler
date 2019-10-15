package sender

import (
	"errors"
	"reflect"
)

var (
	//ErrChanIsClosed 通道已经关闭
	ErrChanIsClosed = errors.New("sender: is closed")
)

//Channel 信道发送者
type Channel struct {
	val reflect.Value
	typ reflect.Type

	elemIsPtr bool
	elemTyp   reflect.Type

	quit chan struct{}
}

//NewChannel 创建信道发送者
func NewChannel(val reflect.Value, typ reflect.Type) Sender {
	ch := &Channel{
		val:  val,
		typ:  typ,
		quit: make(chan struct{}),
	}

	ch.initElemType()
	return ch
}

func (c *Channel) initElemType() {
	c.elemTyp = c.typ.Elem()
	if c.elemTyp.Kind() == reflect.Ptr {
		c.elemIsPtr = true
		c.elemTyp = c.elemTyp.Elem()
	}
}

//Send Sender接口实现
func (c *Channel) Send(val interface{}) error {
	valOf := reflect.ValueOf(val)
	typOf := valOf.Type()
	isPtr := false

	if typOf.Kind() == reflect.Ptr {
		isPtr = true
		typOf = typOf.Elem()
	}

	if typOf != c.elemTyp {
		panic("collector: send different types with the expected " + c.elemTyp.String())
	}

	if isPtr {
		if !c.elemIsPtr {
			valOf = valOf.Elem()
		}
	} else {
		if c.elemIsPtr {
			if valOf.CanAddr() {
				valOf = valOf.Addr()
			} else {
				newVal := reflect.New(c.elemTyp)
				newOf := newVal.Elem()
				newOf.Set(valOf)
				valOf = newVal
			}
		}
	}

	select {
	case <-c.quit:
		return ErrChanIsClosed
	default:
		c.val.Send(valOf)
	}
	return nil
}

//Close 关闭
func (c *Channel) Close() error {
	select {
	case <-c.quit:
		return ErrChanIsClosed
	default:
		close(c.quit)
	}
	c.val.Close()
	return nil
}
