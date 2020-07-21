package gate

import (
	"github.com/LonelySnail/monkey/agent"
	"github.com/LonelySnail/monkey/gate/network"
	"github.com/LonelySnail/monkey/gate/sessions"
	"github.com/LonelySnail/monkey/module"
)

type Gate struct {
	app module.IDefaultApp
	*module.BaseModule
	manage     *sessions.SessionManage
	agentProxy module.ISessionProxy
	options    *Options
	TcpAddr    string
	WSAddr     string
}

func (gt *Gate) OnInit(app module.IDefaultApp, opts ...Option) error {
	gt.app = app
	gt.BaseModule.OnInit(app)
	options := newOptions(opts...)
	gt.options = options
	gt.manage = sessions.NewSessionMange()
	if gt.options.TCPAddr != "" {
		gt.NewTcpServer()
	}

	if gt.options.WsAddr != "" {
		gt.NewWSServer()
	}
	return nil
}

func (gt *Gate) NewTcpServer() {
	server := new(network.TCPServer)
	server.Addr = gt.options.TCPAddr
	server.TLS = gt.options.TLS
	server.CertFile = gt.options.CertFile
	server.KeyFile = gt.options.KeyFile
	server.NewAgent = func(conn *network.TCPConn) network.Agent {
		a := agent.NewAgent(gt.app, conn)
		a.OnInit(gt)
		return a
	}
	server.Start()

}

func (gt *Gate) NewWSServer() {
	server := new(network.WSServer)
	server.Addr = gt.options.WsAddr
	server.TLS = gt.options.TLS
	server.CertFile = gt.options.CertFile
	server.KeyFile = gt.options.KeyFile
	server.NewAgent = func(conn *network.WSConn) network.Agent {
		a := agent.NewAgent(gt.app, conn)
		a.OnInit(gt)
		return a
	}
	server.Start()
}

func (gt *Gate) GetApp() module.IDefaultApp {
	return gt.app
}

func (gt *Gate) Connect(a module.IGateSession) {
	gt.manage.Set(a.GetSessionID(), a)
	if gt.agentProxy != nil {
		gt.agentProxy.Connect(a)
	}
}

func (gt *Gate) DisConnect(a module.IGateSession) {
	gt.manage.Delete(a.GetSessionID())
	if gt.agentProxy != nil {
		gt.agentProxy.DisConnect(a)
	}
}

func (gt *Gate) SetAgentProxy(proxy module.ISessionProxy) {
	gt.agentProxy = proxy
}
