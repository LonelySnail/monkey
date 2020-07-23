package service

import (
	"errors"
	"fmt"
	"github.com/LonelySnail/monkey/codec"
	"github.com/LonelySnail/monkey/logger"
	"github.com/LonelySnail/monkey/module"
	"github.com/LonelySnail/monkey/rpc"
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

func newService(rcv module.Module) (*service, error) {
	service := new(service)
	service.typ = reflect.TypeOf(rcv)
	service.rcv = reflect.ValueOf(rcv)
	name := rcv.GetName()
	if name == "" {
		name = reflect.Indirect(service.rcv).Type().Name() // Type
	}
	if name == "" {
		return nil, errors.New("name is empty")
	}

	if rcv.GetApp() == nil {
		return nil,errors.New("app is nil")
	}
	service.app = rcv.GetApp()
	service.name = name
	service.ch = make(chan []byte,10)
	// Install the methods

	service.method = suitableMethods(service.typ)

	if len(service.method) == 0 {
		//return nil, fmt.Errorf("%s has no methods ", name)
	}
	client, err := rpc.NewRedisClient(fmt.Sprintf("server_test:%s",name), "redis://root@192.168.5.137/6")
	if err != nil {
		return nil, err
	}
	server, err := rpc.NewRpcServer(fmt.Sprintf("server_test:%s",name), "redis://root@192.168.5.137/6",service.ch)
	if err != nil {
		return nil, err
	}
	service.rpcClient = client
	service.rpcServer = server
	go service.handler()
	logger.ZapLog.Error("service:", zap.String("servicePath", service.name))
	return service, nil
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
		// Method needs four ins: receiver, session, *args.
		if mtype.NumIn() != 3 {
			continue
		}
		// First arg must be iface.Isession
		//session := mtype.In(1)
		//if !session.Implements(typeOfSession) {
		//	continue
		//}

		//Second arg need not be a pointer.
		argType := mtype.In(2)
		if !isExportedOrBuiltinType(argType) {
			continue
		}

		// Method needs one out.
		if mtype.NumOut() != 1 {
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType}
	}
	return methods
}

func (s *Service) getService(name string) *service {

	return s.serviceMap[name]
}

//func (s *Service) HandlerRequest(args [][]byte) {
//	defer func() {
//		if r := recover(); r != nil {
//			var rn = ""
//			switch r.(type) {
//
//			case string:
//				rn = r.(string)
//			case error:
//				rn = r.(error).Error()
//			}
//			logger.ZapLog.Error(rn)
//			buff := make([]byte, 1024)
//			runtime.Stack(buff, false)
//			logger.ZapLog.Error(string(buff))
//		}
//	}()
//	if len(args) != 2 {
//		logger.ZapLog.Error("args is error", zap.Any("args", args))
//		return
//	}
//	msg := new(Message)
//	err := codec.GetCodec(s.SerializeType).Decode(args[1], msg)
//	if err != nil {
//		logger.ZapLog.Error(err.Error())
//		return
//	}
//	logger.ZapLog.Info("handlerRequest:", zap.String("servicePath", msg.ServicePath), zap.String("serviceMethod", msg.ServiceMethod), zap.Any("payload", msg.Payload))
//
//	ser := s.getService(msg.ServicePath)
//	if ser == nil {
//		logger.ZapLog.Warn("service is not exist", zap.String("service name", msg.ServicePath))
//		return
//	}
//
//	mty := ser.method[msg.ServiceMethod]
//	if mty == nil {
//		logger.ZapLog.Warn("service method is not exist", zap.String("service method", msg.ServiceMethod))
//		return
//	}
//
//	arg := reflect.ValueOf(msg.Payload)
//	if mty.ArgType.Kind() == reflect.Ptr {
//		arg = reflect.ValueOf(msg.Payload).Elem()
//	}
//	session := new(agent.SessionAgent)
//	err = json.Unmarshal(args[0], session)
//	if err != nil {
//		logger.ZapLog.Warn("session is error", zap.Error(err), zap.String("session", string(args[0])))
//		return
//	}
//
//	if ser.isGo {
//		go ser.call(mty, reflect.ValueOf(session), arg)
//	} else {
//		ser.call(mty, reflect.ValueOf(session), arg)
//	}
//	return
//}

//func (s *Service)Call(path,method string,args ...interface{})  {
//		defer func() {
//			if r := recover(); r != nil {
//				var rn = ""
//				switch r.(type) {
//
//				case string:
//					rn = r.(string)
//				case error:
//					rn = r.(error).Error()
//				}
//				logger.ZapLog.Error(rn)
//				buff := make([]byte, 1024)
//				runtime.Stack(buff, false)
//				logger.ZapLog.Error(string(buff))
//			}
//		}()
//
//	logger.ZapLog.Info("handlerRequest:", zap.String("servicePath", path), zap.String("serviceMethod", method), zap.Any("args", args))
//
//	ser := s.getService(path)
//	if ser == nil {
//		logger.ZapLog.Warn("service is not exist", zap.String("service name", path))
//		return
//	}
//
//	mty := ser.method[method]
//	if mty == nil {
//		logger.ZapLog.Warn("service method is not exist", zap.String("service method", method))
//		return
//	}
//	in := make([]reflect.Value,0)
//	for _,arg := range args {
//		in = append(in,reflect.ValueOf(arg))
//	}
//	if ser.isGo {
//		go ser.call(mty, in)
//	} else {
//		ser.call(mty,in)
//	}
//}

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