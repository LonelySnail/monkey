package app

import (
	"github.com/LonelySnail/monkey/codec"
)

type OptionFn func(opt *Options)

type Options struct {
	version       string
	appType       string
	tcpAddr       string
	wsAddr        string
	SerializeType byte
}

func newOptions(opts ...OptionFn) *Options {
	newOpt := new(Options)
	for _, opt := range opts {
		opt(newOpt)
	}
	if newOpt.SerializeType == 0 {
		newOpt.SerializeType = codec.JSON
	}
	return newOpt
}

func SetAppVersion(ver string) OptionFn {
	return func(opt *Options) { opt.version = ver }
}

func (opts *Options) GetAppVersion() string {
	return opts.version
}

func SetAppType(typ string) OptionFn {
	return func(opt *Options) { opt.appType = typ }
}

func (opts *Options) GetAppType() string {
	return opts.appType
}

func SetTcpAddr(addr string) OptionFn {
	return func(opts *Options) { opts.tcpAddr = addr }
}

func (opts *Options) GetTcpAddr() string {
	return opts.appType
}

func SetWsAddr(addr string) OptionFn {
	return func(opts *Options) { opts.wsAddr = addr }
}

func (opts *Options) GetWSAddr() string {
	return opts.appType
}
func SerializeType(typ byte) OptionFn {
	return func(opt *Options) {
		opt.SerializeType = typ
	}
}

func (opts *Options) GetSerializeType() byte {
	return opts.SerializeType
}
