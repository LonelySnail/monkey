package service

import (
	"errors"
	"fmt"
	"github.com/LonelySnail/monkey/codec"
	"github.com/LonelySnail/monkey/logger"
	"github.com/LonelySnail/monkey/module"
	"go.uber.org/zap"
	"reflect"
	"runtime"
	"sync"
	"unicode"
	"unicode/utf8"
)

// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// Precompute the reflect type for session.
var typeOfSession = reflect.TypeOf((*module.ISession)(nil)).Elem()

type Message struct {
	ServicePath 	string
	ServiceMethod string
	Payload      interface{}
}
type Service struct {
	SerializeType byte
	options      *Options
	serviceMapMu sync.RWMutex        // 保护service提供service记录表的安全(读多写少使用读写锁)
	serviceMap   map[string]*service // server端提供service记录表
}

// 包括service名、服务对应的接收者、接收者类型、service注册的方法以及注册的函数
type service struct {
	isGo    bool
	name     string                   // name of service
	rcv      reflect.Value            // receiver of methods for the service
	typ      reflect.Type             // type of the receiver
	method   map[string]*methodType   // registered methods
	function map[string]*functionType // registered functions
}
//
// 方法类型：包括方法属性、参数属性、响应属性
type methodType struct {
	sync.Mutex // protects counters
	method     reflect.Method
	ArgType    reflect.Type
}

// 函数类型：包括函数属性、参数属性、返回内容属性
type functionType struct {
	sync.Mutex
	fn        reflect.Value
	ArgType   reflect.Type
	ReplyType reflect.Type
}

func NewService(opts ...OptionFn) *Service {
	s := new(Service)
	options := new(Options)
	for _,opt := range opts {
		opt(options)
	}
	s.options = options
	if s.SerializeType ==0 {
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

	service.name = name

	// Install the methods
	service.method = suitableMethods(service.typ)

	if len(service.method) == 0 {
		return nil, fmt.Errorf("%s has no methods ", name)
	}
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
		logger.ZapLog.Fatal(err.Error())
		return err
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
		logger.ZapLog.Fatal(err.Error())
		return err
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
			logger.ZapLog.Info("has wrong number of ins", zap.String("name", method.Name), zap.Int("args num", mtype.NumIn()))
			continue
		}
		// First arg must be iface.Isession
		session := mtype.In(1)
		if !session.Implements(typeOfSession) {
			logger.ZapLog.Info("First arg must be iface.Isession")
			continue
		}

		//Second arg need not be a pointer.
		argType := mtype.In(2)
		if !isExportedOrBuiltinType(argType) {
			logger.ZapLog.Info(mname, zap.Any(" parameter type not exported: ", argType))
			continue
		}

		// Method needs one out.
		if mtype.NumOut() != 1 {
			logger.ZapLog.Info(mname, zap.Int(" has wrong number of outs:", mtype.NumOut()))
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			logger.ZapLog.Info(mname, zap.String(" returns ", returnType.String()))
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType}
	}
	return methods
}

func (s *Service) getService(name string) *service {

	return s.serviceMap[name]
}

func (s *Service) HandlerRequest(session module.ISession,payload []byte) {
	defer func() {
		if r := recover(); r != nil {
			var rn = ""
			switch r.(type) {

			case string:
				rn = r.(string)
			case error:
				rn = r.(error).Error()
			}
			logger.ZapLog.Error(rn)
			buff := make([]byte, 1024)
			runtime.Stack(buff, false)
			logger.ZapLog.Error(string(buff))
		}
	}()
	msg := new(Message)
	err := codec.GetCodec(s.SerializeType).Decode(payload, msg)
	if err != nil {
		logger.ZapLog.Error(err.Error())
		return
	}
	logger.ZapLog.Info("handlerRequest:",zap.String("servicePath",msg.ServicePath),zap.String("serviceMethod",msg.ServiceMethod),zap.Any("payload",msg.Payload))

	ser := s.getService(msg.ServicePath)
	if ser == nil {
		logger.ZapLog.Warn("service is not exist", zap.String("service name", msg.ServicePath))
		return
	}

	mty := ser.method[msg.ServiceMethod]
	if mty == nil {
		logger.ZapLog.Warn("service method is not exist", zap.String("service method", msg.ServiceMethod))
		return
	}

	arg := reflect.ValueOf(msg.Payload)
	if mty.ArgType.Kind() == reflect.Ptr {
		arg =  reflect.ValueOf(msg.Payload).Elem()
	}

	if ser.isGo {
		go ser.call(mty,reflect.ValueOf(session),arg)
	}else {
		ser.call(mty,reflect.ValueOf(session),arg)
	}
	return
}

func (s *service) call(mty *methodType, session, argv reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var rn = ""
			switch r.(type) {

			case string:
				rn = r.(string)
			case error:
				rn = r.(error).Error()
			}
			logger.ZapLog.Error(rn)
			buff := make([]byte, 1024)
			runtime.Stack(buff, false)
			logger.ZapLog.Error(string(buff))
		}
	}()
	function := mty.method.Func
	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{s.rcv, session, argv})
	if len(returnValues) == 0 {
		return
	}
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}

	return nil
}