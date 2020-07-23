package agent

import (
	"bufio"
	"fmt"
	"github.com/LonelySnail/monkey/agent/packet"
	"github.com/LonelySnail/monkey/codec"
	"github.com/LonelySnail/monkey/gate/network"
	"github.com/LonelySnail/monkey/logger"
	"github.com/LonelySnail/monkey/module"
	"sync"

	"runtime"
	"time"
)

type Message struct {
	ServicePath   string
	ServiceMethod string
	Payload      interface{}
}
type Agent struct {
	conn              network.Conn
	r                 *bufio.Reader
	w                 *bufio.Writer
	isClose           bool
	lastHeartbeatTime int64 //上一次发送存储心跳时间
	revNum            int64
	sendNum           int64
	maxPackSize       int
	connTime          time.Time
	writeLock         sync.Mutex
	session           *SessionAgent
	gt                module.IGate
	app               module.IDefaultApp
	code              codec.Codec
}

func NewAgent(app module.IDefaultApp, conn network.Conn) *Agent {
	return &Agent{app: app, conn: conn, maxPackSize: 65535}
}

func (a *Agent) OnInit(gt module.IGate) error {
	a.r = bufio.NewReaderSize(a.conn, 256)
	a.w = bufio.NewWriterSize(a.conn, 256)
	a.isClose = false
	a.revNum = 0
	a.sendNum = 0
	a.lastHeartbeatTime = time.Now().UnixNano()
	a.code = codec.GetCodec(a.app.GetSerializeType())
	a.gt = gt
	return nil
}

func (a *Agent) Run() error {
	defer func() {
		if err := recover(); err != nil {
			buff := make([]byte, 1024)
			runtime.Stack(buff, false)
			logger.ZapLog.Error(fmt.Sprintf("conn.serve() panic(%v)\n info:%s", err, string(buff)))
		}
		a.Close()

	}()
	//go func() {
	//	defer func() {
	//		if err := recover(); err != nil {
	//			buff := make([]byte, 1024)
	//			runtime.Stack(buff, false)
	//			logger.ZapLog.Error(fmt.Sprintf("OverTime panic(%v)\n info:%s", err, string(buff)))
	//		}
	//	}()
	//	select {
	//	case <-time.After(1 * time.Minute):
	//		if a.GetSession() == nil {
	//			//超过一段时间还没有建立mqtt连接则直接关闭网络连接
	//			a.Close()
	//		}
	//	}
	//}()

	//pack,err := packet.UnPacket(a.r)
	//if err != nil || pack.Type != packet.CONNECT {
	//
	//	return err
	//}

	addr := a.conn.RemoteAddr()
	a.session = newSession(addr.String(), addr.Network())
	a.gt.Connect(a)
	a.listenAndLoop()
	return nil
}

func (a *Agent) listenAndLoop() {
	defer func() {
		if err := recover(); err != nil {
			buff := make([]byte, 1024)
			runtime.Stack(buff, false)
			logger.ZapLog.Error(fmt.Sprintf("conn.serve() panic(%v)\n info:%s", err, string(buff)))
		}
		a.Close()

	}()
	go a.Flusher()
	a.readMsg()
}
func (a *Agent) isClosed() bool {
	return a.isClose
}

//  给客户端返回信息
func (a *Agent) Flusher() {
	for !a.isClosed() {
		a.writeLock.Lock()
		if a.isClosed() {
			a.writeLock.Unlock()
			break
		}
		if a.w.Buffered() > 0 {
			if err := a.w.Flush(); err != nil {
				a.writeLock.Unlock()
				break
			}
		}
		a.writeLock.Unlock()
	}
}

func (a *Agent) readMsg() {
	for !a.isClose {
		a.conn.SetDeadline(time.Now().Add(time.Second * 90))
		p, err := packet.UnPacket(a.r)
		if err != nil {
			break
		}
		a.handlerMsg(p)
	}

}

func (a *Agent) handlerMsg(pack *packet.Pack) {
	switch pack.Type {
	case packet.CONNECT:

	case packet.PING:
		a.lastHeartbeatTime = time.Now().Unix()
	case packet.DATA:
		msg := new(Message)
		err :=a.code.Decode(pack.Payload,msg)
		if err != nil {
			return
		}
		a.app.CallNR(msg.ServicePath,msg.ServiceMethod,a.session,msg.Payload)
	case packet.DISCONNECT:
	}
}

func (a *Agent) Close() error {
	if a.conn != nil {
		a.conn.Close()
	}
	a.isClose = true
	return nil
}

func (a *Agent) OnClose() error {
	defer func() {
		if err := recover(); err != nil {
			buff := make([]byte, 1024)
			runtime.Stack(buff, false)
			logger.ZapLog.Error(fmt.Sprintf("agent OnClose panic(%v)\n info:%s", err, string(buff)))
		}
	}()
	a.isClose = true
	a.gt.DisConnect(a) //发送连接断开的事件
	return nil
}

func (a *Agent) GetSessionID() string {
	return a.session.GetSessionID()
}

func (a *Agent) GetIP() string {
	return a.session.IP
}

func (a *Agent) GetNetWork() string {
	return a.session.Network
}

func (a *Agent) SendMsg(p []byte) (int, error) {
	return a.w.Write(p)
}

func (a *Agent) Send(p []byte) {
	a.app.Call("gate", "Send", a.GetSessionID(),p)
}
