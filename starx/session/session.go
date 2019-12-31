package session

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/chrislonng/starx/log"
	"github.com/chrislonng/starx/service"
)

type NetworkEntity interface {
	ID() int64
	Send([]byte) error
	Push(session *Session, route string, v interface{}) error
	Response(session *Session, v interface{}) error
	Call(session *Session, route string, reply interface{}, args ...interface{}) error
	Close()
}

var (
	ErrIllegalUID       = errors.New("illegal uid")
	ErrKeyNotFound      = errors.New("current session does not contain key")
	ErrWrongValueType   = errors.New("current key has different data type")
	ErrReplyShouldBePtr = errors.New("reply should be a pointer")
)

// This session type as argument pass to Handler method, is a proxy session
// for frontend session in frontend server or backend session in backend
// server, correspond frontend session or backend session id as a field
// will be store in type instance
//
// This is user sessions, does not contain raw sockets information
type Session struct {
	ID        int64                  // session global unique id
	Uid       int64                  // binding user id
	Entity    NetworkEntity          // raw session id, agent in frontend server, or acceptor in backend server
	LastID    uint                   // last request id
	data      map[string]interface{} // session data store
	lastTime  int64                  // last heartbeat time
	serverIDs map[string]string      // map of server type -> server id
}

// Create new session instance
func New(entity NetworkEntity) *Session {
	return &Session{
		ID:        service.Connections.SessionID(),
		Entity:    entity,
		data:      make(map[string]interface{}),
		lastTime:  time.Now().Unix(),
		serverIDs: make(map[string]string),
	}
}

func (s *Session) ServerID(svrType string) string {
	id, ok := s.serverIDs[svrType]
	if !ok {
		return ""
	}
	return id
}

// Set server id of the special type, delete type when id empty
func (s *Session) SetServerID(svrType, svrID string) {
	svrType = strings.TrimSpace(svrType)
	svrID = strings.TrimSpace(svrID)

	if svrType == "" {
		log.Errorf("empty server type")
		return
	}

	if svrID == "" {
		delete(s.serverIDs, svrType)
		return
	}
	s.serverIDs[svrType] = svrID
}

// Session send packet data
func (s *Session) Send(data []byte) error {
	return s.Entity.Send(data)
}

// Push message to session
func (s *Session) Push(route string, v interface{}) error {
	return s.Entity.Push(s, route, v)
}

// Response message to session
func (s *Session) Response(v interface{}) error {
	return s.Entity.Response(s, v)
}

func (s *Session) Bind(uid int64) error {
	if uid < 1 {
		log.Errorf("uid invalid: %d", uid)
		return ErrIllegalUID
	}
	s.Uid = uid
	return nil
}

func (s *Session) Call(route string, reply interface{}, args ...interface{}) error {
	if reflect.TypeOf(reply).Kind() != reflect.Ptr {
		return ErrReplyShouldBePtr
	}
	return s.Entity.Call(s, route, reply, args...)
}

func (s *Session) Close() {
	s.Entity.Close()
}

func (s *Session) Remove(key string) {
	delete(s.data, key)
}

func (s *Session) Set(key string, value interface{}) {
	s.data[key] = value
}

func (s *Session) HasKey(key string) bool {
	_, has := s.data[key]
	return has
}

func (s *Session) Int(key string) int {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int8(key string) int8 {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int8)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int16(key string) int16 {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int16)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int32(key string) int32 {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int32)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int64(key string) int64 {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int64)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint(key string) uint {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint8(key string) uint8 {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint8)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint16(key string) uint16 {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint16)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint32(key string) uint32 {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint32)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint64(key string) uint64 {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint64)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Float32(key string) float32 {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(float32)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Float64(key string) float64 {
	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(float64)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) String(key string) string {
	v, ok := s.data[key]
	if !ok {
		return ""
	}

	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}

func (s *Session) Value(key string) interface{} {
	return s.data[key]
}

// Retrieve all session state
func (s *Session) State() map[string]interface{} {
	return s.data
}

// Restore session state after reconnect
func (s *Session) Restore(data map[string]interface{}) {
	s.data = data
}

func (s *Session) Clear() {
	log.Debugf("Clear session data: Id=%d, Uid=%d", s.ID, s.Uid)
	s.data = map[string]interface{}{}
}
