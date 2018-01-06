package main

import (
	"sync"
	"time"

	"github.com/emsihyo/bi"
	nats "github.com/nats-io/go-nats"
)

//Gate Gate
type Gate struct {
	bi.BI
	mut     sync.Mutex
	clients map[string]*Session
}

//NewGate NewGate
func NewGate() *Gate {
	gate := &Gate{clients: map[string]*Session{}}
	return gate
}

func (gate *Gate) init() {
	nc, err := nats.Connect(nats.DefaultURL)
	if nil != err {
		panic(err)
	}
	gate.On(bi.Connection, func(sess *Session) {
		gate.mut.Lock()
		gate.clients[sess.GetID()] = sess
		gate.mut.Unlock()
	})
	gate.On(bi.Disconnection, func(sess *Session) {
		gate.mut.Lock()
		delete(gate.clients, sess.GetID())
		gate.mut.Unlock()
	})
	gate.On("publish", func(sess *Session, pub *Publish) {
		nc, err := nats.Connect(nats.DefaultURL)
		if nil == err {
			nc.Publish(pub.Subject, pub.Data)
			nc.Close()
		}
	})
	gate.On("request", func(sess *Session, req *Request) *Response {
		nc, err := nats.Connect(nats.DefaultURL)
		var data []byte
		if nil == err {
			var msg *nats.Msg
			msg, err = nc.Request(req.Subject, req.Data, time.Second*10)
			nc.Close()
			if nil == err {
				data = msg.Data
			}
		}
		resp := &Response{Err: err.Error(), Data: data}
		return resp
	})
}
