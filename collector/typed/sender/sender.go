package sender

//Sender 发送者接口
type Sender interface {
	Send(val interface{}) error
	Close() error
}
