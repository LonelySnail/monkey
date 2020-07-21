package service

import (
	"github.com/LonelySnail/monkey/logger"
	"github.com/LonelySnail/monkey/module"
	"go.uber.org/zap"
	"reflect"
	"runtime"
	"sync"
)

type Message struct {
	ServicePath   string
	ServiceMethod string
	Payload       interface{}
}

// 包括service名、服务对应的接收者、接收者类型、service注册的方法以及注册的函数
type service struct {
	app  module.IDefaultApp
	isGo      bool
	rpcClient IRpcClient
	rpcServer IRpcServer
	name      string                   // name of service
	rcv       reflect.Value            // receiver of methods for the service
	typ       reflect.Type             // type of the receiver
	method    map[string]*methodType   // registered methods
	function  map[string]*functionType // registered functions
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

func (s *service) Call(msg *Message) {
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
	mty := s.method[msg.ServiceMethod]
	if mty == nil {
		logger.ZapLog.Warn("service method is not exist", zap.String("service method", msg.ServiceMethod))
		return
	}

	//arg := reflect.ValueOf(msg.Payload)
	//if mty.ArgType.Kind() == reflect.Ptr {
	//	arg =  reflect.ValueOf(msg.Payload).Elem()
	//}

	//s.call(mty,reflect.ValueOf(session),arg)
}

func (s *service) call(mty *methodType, args []reflect.Value) (err error) {
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
	returnValues := function.Call(args)
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
