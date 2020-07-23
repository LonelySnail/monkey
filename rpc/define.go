package rpc

import (
	"context"
	"time"
)

type IRpcClient interface {
	Call(msg *RpcMsg) error
	CallNR(msg *RpcMsg) error
}

type IRpcServer interface {

}

type RpcMsg struct {
	ID         string
	Method     string
	Reply      bool
	Args       [][]byte
	ArgsType[]string
}

func getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

//func NewRpc(typ string,rpcFunc rpcFunction)  *Rpc{
//	var r *Rpc
//	var err error
//	switch typ {
//	case "redis":
//		r,err = NewRedisRpc()
//	default:
//		err = errors.New("no rpc")
//	}
//	if err != nil {
//		logger.ZapLog.Fatal(err.Error())
//	}
//	r.rpcFunc = rpcFunc
//	r.callChan = make(chan *CallMsg,100)
//	go r.CallHandler()
//	return r
//}

//func NewRedisRpc() (*Rpc,error) {
//	r := new(Rpc)
//	c,err := NewRedisClient("test","redis://root@192.168.5.137/6")
//	if err != nil {
//		return  r,err
//	}
//	s,err := NewRedisServer("test","redis://root@192.168.5.137/6")
//	if err != nil {
//		return r,err
//	}
//	r.Client = c
//	r.Server =s
//
//	return  r,nil
//}

////  需要等待返回值
//func (r *Rpc) Call(session iface.ISession, msg *codec.Message) (err error) {
//	ctx := getContext()
//	ch := make(chan int, 1)
//	call, err := r.Go(ctx, session, msg, ch)
//	select {
//	case <-ctx.Done():
//		r.pending.Delete(call.Seq)
//		err = ctx.Err()
//	case <-ch:
//
//	}
//	return err
//}
//
//// 异步调用
//func (r *Rpc) Go(ctx context.Context, session iface.ISession, msg *codec.Message, ch chan int) (call *CallMsg, err error) {
//	call = r.newCallMsg(session, msg)
//	r.pending.Store(call.Seq, call)
//	err = r.Client.Go(call)
//	ch <- 1
//	return call, err
//}
//
//func (r *Rpc) CallHandler() {
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
//			buf := make([]byte, 1024)
//			l := runtime.Stack(buf, false)
//			errstr := string(buf[:l])
//			logger.ZapLog.Error(rn, zap.String("Stack", errstr))
//		}
//	}()
//	go r.Server.requestHandler(r.callChan)
//
//	select {
//	case call, ok := <-r.callChan:
//		if ok {
//			r.rpcFunc(call)
//		}
//
//	}
//}
