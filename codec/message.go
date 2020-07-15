package codec

import "sync"

type Message struct {
	ServicePath   string
	ServiceMethod string
	Payload       []byte
}

var msgPool = sync.Pool{
	New: func() interface{} {
		return new(Message)
	},
}

func NewMessage(path, method string, payload []byte) *Message {
	return &Message{
		ServicePath:   path,
		ServiceMethod: method,
		Payload:       payload,
	}
}

func (msg *Message) Reset() {
	msg.ServiceMethod = ""
	msg.ServicePath = ""
	msg.Payload = msg.Payload[:0]
}

func (msg *Message) SetMsg(path, method string, payload []byte) {
	msg.ServicePath = path
	msg.ServiceMethod = method
	msg.Payload = payload
}

func Get() *Message {
	msg, _ := msgPool.Get().(*Message)
	return msg
}

func (msg *Message) Put() {
	msg.Reset()
	msgPool.Put(msg)
}
