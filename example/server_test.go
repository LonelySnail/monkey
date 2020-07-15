package example

import (
	"encoding/json"
	"fmt"
	"github.com/LonelySnail/monkey"
	"github.com/LonelySnail/monkey/app"
	"github.com/LonelySnail/monkey/gate"
	"github.com/LonelySnail/monkey/module"
	"testing"
)

type mgate struct {
	*gate.Gate
}

type Login struct {
	app module.IDefaultApp
}
func newGate() *mgate {
	g := new(mgate)
	g.Gate = new(gate.Gate)
	return g
}

func (m *mgate)GetName()string  {
	return "gate"
}

func (m *mgate)GetType()string  {
	return "2"
}

func (m *mgate) OnInit(app module.IDefaultApp)error {
	m.Gate.OnInit(app,gate.TCPAddr(":3598"))
	return nil
}

func(m *mgate)Login(session module.ISession,name string) (err error) {
	fmt.Println(name,"666666")
	return
}

func (l *Login) GetName()string  {
	return "login"
}

func (l *Login) GetType()string  {
	return "login"
}

func (l *Login) OnInit(app module.IDefaultApp)error  {
	l.app = app
	return nil
}
func (l *Login) GetApp()module.IDefaultApp  {
	return l.app
}

func (l *Login)Login(session module.ISession,arg map[string]interface{})(err error)  {
	fmt.Println(arg,"666666")
	a,_ :=json.Marshal("hello world")
	session.SendMsg(a)
	return
}
func TestServer(t *testing.T)  {
	a := monkey.NewDefaultApp(app.SetTcpAddr(":3598"))
	a.Start(newGate(),new(Login))
}
