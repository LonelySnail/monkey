package sessions

import (
	"github.com/LonelySnail/monkey/module"
	"sync"
)

type SessionManage struct {
	lock       sync.RWMutex
	sessionMap map[string]module.IGateSession //  session map
	num        int                            // session 数量
}

func NewSessionMange() *SessionManage {
	return &SessionManage{sessionMap: make(map[string]module.IGateSession)}
}

func (s *SessionManage) GetSession(key string) (module.IGateSession, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	session, ok := s.sessionMap[key]
	return session, ok
}

func (s *SessionManage) Set(key string, session module.IGateSession) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sessionMap[key] = session
	s.num++
}

func (s *SessionManage) Delete(key string) {
	s.lock.Lock()
	s.lock.Unlock()
	delete(s.sessionMap, key)
	s.num--
}
