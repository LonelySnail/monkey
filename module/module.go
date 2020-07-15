package module

type IDefaultApp interface {
	Start(mods ...Module)
	OnInit()
	OnStop()
	GetTcpAddr()string
	GetWSAddr() string
	GetService() IService
}

type Module interface {
	GetName()	string
	GetType()   string
	OnInit(app IDefaultApp)    error
	GetApp() IDefaultApp
}

type ISession interface {
	SendMsg(msg []byte) (int,error)
}

type IService interface {
	HandlerRequest(session ISession,payload []byte)
}