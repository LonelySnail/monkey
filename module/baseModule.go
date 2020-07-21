package module

type BaseModule struct {
	app  IDefaultApp
}

func (base *BaseModule) OnInit(app IDefaultApp) {
	base.app = app
}

func (base *BaseModule) GetName() string {
	return ""
}

func (base *BaseModule) GetType() string {
	return ""
}

func (base *BaseModule) GetApp() IDefaultApp {
	return base.app
}

