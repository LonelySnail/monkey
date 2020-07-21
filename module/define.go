package module

type IDefaultApp interface {
	Start(mods ...Module)
	OnInit()
	OnStop()
	GetTcpAddr() string
	GetWSAddr() string
	RequestCall(args [][]byte)
	GetSerializeType()byte
	Call(path, method string, args ...interface{})
}

type Module interface {
	GetName() string
	GetType() string
	OnInit(app IDefaultApp) error
	GetApp() IDefaultApp
}

//type IAgent interface {
//	GetSessionID()    string
//	GetSession()     IGateSession
//}

type IGate interface {
	Connect(session IGateSession)
	DisConnect(session IGateSession)
}

type IGateSession interface {
	GetSessionID() string
	GetIP() string
	GetNetWork() string
	SendMsg(payload []byte) (int, error)
}

//  session 代理
type ISessionProxy interface {
	Connect(session IGateSession)
	DisConnect(session IGateSession)
}
