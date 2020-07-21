package agent

import (
	utils "github.com/LonelySnail/monkey/util"
)

type SessionAgent struct {
	ID        string            `json:"id,omitempty"`
	IP        string            `json:"IP,omitempty"`
	Network   string            `json:"Network,omitempty"`
	UserId    string            `json:"UserId,omitempty"`
	SessionId string            `json:"SessionId,omitempty"`
	Settings  map[string]string `json:"Settings,omitempty"`
}

func newSession(ip, network string) *SessionAgent {
	agent := new(SessionAgent)
	agent.ID = utils.UUid()
	agent.IP = ip
	agent.Network = network
	return agent
}

func (ses *SessionAgent) Clone() *SessionAgent {
	return ses
}

func (ses *SessionAgent) GetSessionID() string {
	return ses.ID
}
