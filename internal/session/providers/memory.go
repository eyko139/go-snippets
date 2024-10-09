// NOTE: this is only to be used if mongo session db is unavailable
package providers

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/eyko139/go-snippets/internal/session"
)

// The Provider here is just a list in memory
// NOTE: Struct properties that are not explicitly initialized are set to their zero value.
// For mutex, that means an unlocked mutex!
var memSessionPder = &InMemorySessionProvider{list: list.New()}

type InMemorySessionProvider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

type SessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[string]interface{}
}

func (st *SessionStore) Set(key string, value interface{}) error {
	st.value[key] = value
	memSessionPder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) Get(key string) interface{} {
	memSessionPder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	}
	return nil
}

func (st *SessionStore) Delete(key string) error {
	delete(st.value, key)
	memSessionPder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

func (pder *InMemorySessionProvider) SessionInit(sid string) (session.Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[string]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess, nil
}

func (pder *InMemorySessionProvider) SessionRead(sid string) (session.Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	}
	sess, err := pder.SessionInit(sid)
	return sess, err
}

func (pder *InMemorySessionProvider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *InMemorySessionProvider) SessionGC(maxLifeTime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}

		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (pder *InMemorySessionProvider) SessionUpdate(sid string) error {
	fmt.Println(pder)
	pder.lock.Lock()
	defer pder.lock.Unlock()

	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
	}
	return nil
}

func InitMemorySession() {
	memSessionPder.sessions = make(map[string]*list.Element)
	session.Register("memory", memSessionPder)
}
