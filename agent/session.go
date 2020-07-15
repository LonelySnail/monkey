package agent

import utils "github.com/LonelySnail/monkey/util"

type ISession interface {

}

type Session struct {
	ID					 string				`json:"id,omitempty"`
	IP                   string            `json:"IP,omitempty"`
	Network              string            `json:"Network,omitempty"`
	UserId               string            `json:"UserId,omitempty"`
	SessionId            string            `json:"SessionId,omitempty"`
	Settings             map[string]string `json:"Settings,omitempty"`
}

func newSession(ip,network string)  *Session{
	session := new(Session)
	session.ID = utils.UUid()
	session.IP = ip
	session.Network = network
	return session
}