package service

import (
	"github.com/LonelySnail/monkey/codec"
	"github.com/LonelySnail/monkey/logger"
	"github.com/LonelySnail/monkey/module"
	"go.uber.org/zap"
	"reflect"
	"sync"
	"unicode"
	"unicode/utf8"
)

// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// Precompute the reflect type for session.
//var typeOfSession = reflect.TypeOf(()(nil)).Elem()

type Service struct {
	SerializeType byte
	options       *Options
	serviceMapMu  sync.RWMutex        // 保护service提供service记录表的安全(读多写少使用读写锁)
	serviceMap    map[string]*service // server端提供service记录表
}

func NewService(opts ...OptionFn) *Service {
	s := new(Service)
	options := new(Options)
	for _, opt := range opts {
		opt(options)
	}
	s.options = options
	if s.SerializeType == 0 {
		s.SerializeType = codec.JSON
	}

	s.serviceMap = make(map[string]*service)
	return s
}

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func (s *Service) Register(rcv module.Module) error {
	s.serviceMapMu.Lock()
	defer s.serviceMapMu.Unlock()
	if s.serviceMap == nil {
		s.serviceMap = make(map[string]*service)
	}

	service, err := newService(rcv)
	if err != nil {
		//logger.ZapLog.Fatal(err.Error())
		//return err
	}
	service.isGo = false
	s.serviceMap[service.name] = service
	return nil
}

func (s *Service) RegisterGo(rcv module.Module) error {
	s.serviceMapMu.Lock()
	defer s.serviceMapMu.Unlock()
	if s.serviceMap == nil {
		s.serviceMap = make(map[string]*service)
	}

	service, err := newService(rcv)
	if err != nil {
		//logger.ZapLog.Fatal(err.Error())
		//return err
	}
	service.isGo = true
	s.serviceMap[service.name] = service
	return nil
}

func suitableMethods(typ reflect.Type) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs four ins: receiver, session, args.
		//if mtype.NumIn() != 3 {
		//	continue
		//}
		// First arg must be iface.Isession
		//session := mtype.In(1)
		//if !session.Implements(typeOfSession) {
		//	continue
		//}

		//Second arg need not be a pointer.
		//argType := mtype.In(2)
		//if !isExportedOrBuiltinType(argType) {
		//	continue
		//}

		// Method needs one out.
		if mtype.NumOut() != 1 {
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			continue
		}
		methods[mname] = &methodType{method: method}
	}
	return methods
}

func (s *Service) getService(name string) *service {

	return s.serviceMap[name]
}

func(s *Service)CallNR (path, method string, args ...interface{})  {
	ser := s.getService(path)
	if ser == nil {
		logger.ZapLog.Warn("service is not exist", zap.String("service name", path))
		return
	}

	argsType := make([]string,len(args))
	params := make([][]byte,len(args))
	var err error
	for i,arg := range args {
		argsType[i],params[i],err = ArgsTypeAnd2Bytes(arg)
		if err != nil {
			return
		}
	}

	ser.CallNR(method,argsType,params)
}