package service

import (
	"encoding/json"
	"fmt"
	"github.com/LonelySnail/monkey/logger"
	"github.com/LonelySnail/monkey/module"
	"github.com/LonelySnail/monkey/rpc"
	utils "github.com/LonelySnail/monkey/util"
	"go.uber.org/zap"
	"errors"
	"runtime"

	//utils "github.com/LonelySnail/monkey/util"
	"reflect"
	"sync"
)

// 包括service名、服务对应的接收者、接收者类型、service注册的方法以及注册的函数
type service struct {
	app  module.IDefaultApp
	isGo      bool
	rpcClient rpc.IRpcClient
	rpcServer rpc.IRpcServer
	name      string                   // name of service
	rcv       reflect.Value            // receiver of methods for the service
	typ       reflect.Type             // type of the receiver
	method    map[string]*methodType   // registered methods
	function  map[string]*functionType // registered functions
	ch   chan  []byte
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
	service.method = suitableMethods(service.typ)
	if len(service.method) == 0 {
		return nil, fmt.Errorf("%s has no methods ", name)
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


func (s *service) CallNR(method string,argsType []string,args [][]byte)  {
	msg := &rpc.RpcMsg{
		ID: utils.UUid(),
		Method: method,
		Reply: false,
		Args: args,
		ArgsType: argsType,
	}

	s.rpcClient.CallNR(msg)
}

func (s *service)handler()  {
	for body := range s.ch {
		msg := new(rpc.RpcMsg)
		err := json.Unmarshal(body,msg)
		if err != nil {
			continue
		}
		err = s.call(msg)
	}
}

func (s *service)call(msg *rpc.RpcMsg) error {
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

	mty := s.method[msg.Method]
	if mty == nil {
		return fmt.Errorf("service [%s] method is not exist",msg.Method)
	}

	args := msg.Args
	fn := mty.method.Func
	in := make([]reflect.Value,len(args)+1)
	in[0]=s.rcv

	for i,typ := range msg.ArgsType {
		ty,err := Bytes2Args(s.app,typ,args[i])
		if err != nil {
			return err
		}
		inx := i+1
		switch v2 := ty.(type) { //多选语句switch
		case nil:
			in[inx] = reflect.Zero(fn.Type().In(i))
		case []uint8:
			if reflect.TypeOf(ty).AssignableTo(fn.Type().In(i)) {
				//如果ty "继承" 于接受参数类型
				in[inx] = reflect.ValueOf(ty)
			} else {
				elemp := reflect.New(fn.Type().In(i))
				err := json.Unmarshal(v2, elemp.Interface())
				if err != nil {
					in[inx] = reflect.ValueOf(ty)
				} else {
					in[inx] = elemp.Elem()
				}
			}
		default:
			in[inx] = reflect.ValueOf(ty)
		}
	}

	// Invoke the method, providing a new value for the reply.
	returnValues := fn.Call(in)
	if len(returnValues) == 0 {
		return fmt.Errorf("no returnValues")
	}
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}
	return nil
}