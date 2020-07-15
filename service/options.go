package service

type OptionFn func(opt *Options)

type Options struct {
	SerializeType byte
}

func SerializeType(typ byte) OptionFn {
	return func(opt *Options) {
		opt.SerializeType = typ
	}
}

func (opts *Options) GetSerializeType() byte {
	return opts.SerializeType
}