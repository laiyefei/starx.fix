package cluster

import (
	"encoding/json"

	"github.com/chrislonng/starx/log"
	"github.com/chrislonng/starx/session"
)

type Manager struct {
	Name    string
	Counter int
}

// Component interface methods
func (m *Manager) Init() {
	m.Name = "ManagerComponenet"
	log.Infof("manager component initialized")
}
func (*Manager) AfterInit()      {}
func (*Manager) BeforeShutdown() {}
func (*Manager) Shutdown()       {}

// attachment methods
func (m *Manager) UpdateServer(session *session.Session, data []byte) error {
	var newServerInfo *ServerConfig
	err := json.Unmarshal(data, newServerInfo)
	if err != nil {
		return err
	}
	UpdateServer(newServerInfo)
	return nil
}

func (m *Manager) RegisterServer(session *session.Session, data []byte) error {
	var newServerInfo *ServerConfig
	err := json.Unmarshal(data, newServerInfo)
	if err != nil {
		return err
	}
	log.Infof("new server connected in")
	Register(newServerInfo)
	return nil
}

func (m *Manager) RemoveServer(session *session.Session, data []byte) error {
	var srvId string
	err := json.Unmarshal(data, &srvId)
	if err != nil {
		return err
	}
	RemoveServer(srvId)
	return nil
}
