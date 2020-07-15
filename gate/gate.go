package gate

import (
	"github.com/LonelySnail/monkey/agent"
	"github.com/LonelySnail/monkey/gate/network"
	"github.com/LonelySnail/monkey/module"
)

type Gate struct {
	app  module.IDefaultApp
	options  *Options
	TcpAddr   string
	WSAddr    string
}

func (gt *Gate)OnInit(app module.IDefaultApp,opts ...Option) error {
	gt.app = app
	options := newOptions(opts...)
	gt.options = options
	if gt.options.TCPAddr != "" {
		gt.NewTcpServer()
	}

	if gt.options.WsAddr != "" {
		gt.NewWSServer()
	}
	return nil
}

func  (gt *Gate)NewTcpServer()  {
	server := new(network.TCPServer)
	server.Addr = gt.options.TCPAddr
	server.Addr = gt.options.TCPAddr
	server.TLS = gt.options.TLS
	server.CertFile = gt.options.CertFile
	server.KeyFile = gt.options.KeyFile
	server.NewAgent = func(conn *network.TCPConn) network.Agent {
		a := agent.NewAgent(gt.app,conn)
		a.OnInit()
		return a
	}
	server.Start()

}

func  (gt *Gate)NewWSServer()  {
	server := new(network.WSServer)
	server.Addr = gt.options.TCPAddr
	server.Addr = gt.options.TCPAddr
	server.TLS = gt.options.TLS
	server.CertFile = gt.options.CertFile
	server.KeyFile = gt.options.KeyFile
	server.NewAgent = func(conn *network.WSConn) network.Agent {
		a := agent.NewAgent(gt.app,conn)
		a.OnInit()
		return a
	}
	server.Start()
}

func (gt *Gate)GetApp() module.IDefaultApp {
	return gt.app
}