package app

import (
	"flag"
	"fmt"
	"github.com/LonelySnail/monkey/logger"
	"github.com/LonelySnail/monkey/module"
	"github.com/LonelySnail/monkey/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type DefaultApp struct {
	service *service.Service
	options *Options
	mods    []module.Module
}

func NewDefaultApp(opts ...OptionFn) *DefaultApp {
	app := new(DefaultApp)
	options := new(Options)
	for _, opt := range opts {
		opt(options)
	}

	app.options = options
	app.service = service.NewService(service.SerializeType(app.options.GetSerializeType()))
	return app
}

func (app *DefaultApp) Start(mods ...module.Module) {
	typ := flag.String("typ", "", "server type")
	flag.Parse()
	err := app.runMods(*typ, mods...)
	if err != nil {
		panic(err)
	}
	app.OnInit()
}

func (app *DefaultApp) runMods(typ string, mods ...module.Module) error {
	for _, mod := range mods {
		app.mods = append(app.mods, mod)
		//if typ == "" || mod.GetType() == typ {
		//	err := mod.OnInit(app)
		//	if err != nil {
		//		return fmt.Errorf("err:[%w],name:%s", err,mod.GetName())
		//	}
		//	app.mods = append(app.mods,mod)
		//}
	}
	return nil
}

func (app *DefaultApp) OnInit() {
	for _, mod := range app.mods {
		mod.OnInit(app)

	}
	app.RegisterGo()
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	sig := <-c

	//如果一分钟都关不了则强制关闭
	timeout := time.NewTimer(1 * time.Minute)
	wait := make(chan struct{})
	go func() {
		app.OnStop()
		wait <- struct{}{}
	}()
	select {
	case <-timeout.C:
		panic(fmt.Sprintf("monkey close timeout (signal: %v)", sig))
	case <-wait:
		logger.ZapLog.Info(fmt.Sprintf("monkey closing down (signal: %v)", sig))
	}
	return
}

func (app *DefaultApp) GetService() *service.Service {
	return app.service
}

func (app *DefaultApp) Register() {
	for _, server := range app.mods {
		err := app.service.Register(server)
		if err != nil {
			panic(err)
		}
	}
}

func (app *DefaultApp) RegisterGo() {
	for _, server := range app.mods {
		err := app.service.RegisterGo(server)
		if err != nil {
			panic(err)
		}
	}
}
func (app *DefaultApp) GetTcpAddr() string {
	return app.options.tcpAddr
}

func (app *DefaultApp) GetWSAddr() string {
	return app.options.wsAddr
}

func (app *DefaultApp) OnStop() {

}

func (app *DefaultApp) RequestCall(args [][]byte) {
	//app.service.HandlerRequest(args)
}

func (app *DefaultApp) RequestCallGetSerializeType()byte{
	return app.options.SerializeType
}

func (app *DefaultApp) Invoke(path, method string, args ...interface{}) {

}

func (app *DefaultApp) InvokeNR(path, method string, args ...interface{}) {

}

func (app *DefaultApp) Call(path, method string, args ...interface{}) {

}

func (app *DefaultApp) CallNR(path, method string, args ...interface{}) {

}

