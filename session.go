package main

import (
	"sync"
	"time"

	"github.com/emsihyo/bi"
)

//Session Session
type Session struct {
	bi.Session
	UserID string
	topics map[string]string
	mut    sync.Mutex
}

//NewSession NewSession
func NewSession(conn bi.Conn, protocol bi.Protocol, timeout time.Duration) *Session {
	return &Session{Session: *bi.NewSession(conn, protocol, timeout), topics: map[string]string{}}
}

func (sess *Session) didSubscribeTopic(topic string) {
	sess.mut.Lock()
	sess.topics[topic] = topic
	sess.mut.Unlock()
}
func (sess *Session) didUnsubscribeTopic(topic string) {
	sess.mut.Lock()
	delete(sess.topics, topic)
	sess.mut.Unlock()
}
func (sess *Session) isSubscribedTopic(topic string) bool {
	sess.mut.Lock()
	defer sess.mut.Unlock()
	_, ok := sess.topics[topic]
	return ok
}
