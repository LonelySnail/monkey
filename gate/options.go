//Package gate 网关配置
package gate

import (
	"time"
)

//Option 网关配置项
type Option func(*Options)

//Options 网关配置项
type Options struct {
	BufSize         int
	MaxPackSize     int
	TLS             bool
	TCPAddr         string
	WsAddr          string
	CertFile        string
	KeyFile         string
	Heartbeat       time.Duration
	OverTime        time.Duration
	AgentLearner
}

type AgentLearner interface {
	Connect()
	DisConnect()
}

//NewOptions 网关配置项
func newOptions(opts ...Option) *Options {
	opt := &Options{
		BufSize:         2048,
		MaxPackSize:     65535,
		Heartbeat:       time.Minute,
		OverTime:        time.Second * 10,
		TLS:             false,
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

//BufSize 单个连接网络数据缓存大小
func BufSize(s int) Option {
	return func(o *Options) {
		o.BufSize = s
	}
}

//MaxPackSize 单个协议包数据最大值
func MaxPackSize(s int) Option {
	return func(o *Options) {
		o.MaxPackSize = s
	}
}

//Heartbeat 心跳时间
func Heartbeat(s time.Duration) Option {
	return func(o *Options) {
		o.Heartbeat = s
	}
}

//OverTime 超时时间
func OverTime(s time.Duration) Option {
	return func(o *Options) {
		o.OverTime = s
	}
}

//Tls Tls
// Deprecated: 因为命名规范问题函数将废弃,请用TLS代替
func Tls(s bool) Option {
	return func(o *Options) {
		o.TLS = s
	}
}

//TLS TLS
func TLS(s bool) Option {
	return func(o *Options) {
		o.TLS = s
	}
}

// TCPAddr tcp监听端口
func TCPAddr(s string) Option {
	return func(o *Options) {
		o.TCPAddr = s
	}
}

// WsAddr websocket监听端口
func WsAddr(s string) Option {
	return func(o *Options) {
		o.WsAddr = s
	}
}

// CertFile TLS 证书cert文件
func CertFile(s string) Option {
	return func(o *Options) {
		o.CertFile = s
	}
}

// KeyFile TLS 证书key文件
func KeyFile(s string) Option {
	return func(o *Options) {
		o.KeyFile = s
	}
}

func SetAgentLearner(learner AgentLearner) Option {
	return func(o *Options) {
		o.AgentLearner = learner
	}
}