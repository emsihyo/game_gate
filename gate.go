package main

import (
	"time"

	"github.com/emsihyo/bi"
	proto "github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/go-nats"
)

//Gate Gate
type Gate struct {
	bi.BI
	id  string
	ns  *Nats
	sub *nats.Subscription
}

//NewGate NewGate
func NewGate() *Gate {
	gate := &Gate{ns: NewNats(nats.DefaultURL, 100, &bi.ProtobufProtocol{})}
	return gate
}

func (gate *Gate) init() {
	id := gate.id
	idx := len(id) + 1
	var err error
	gate.sub, err = gate.ns.Subscribe3(id+".*", idx, func(method string, v interface{}) {
		switch method {
		case ""
		}
	})
	if nil != err {
		panic(err)
	}
	gate.On(bi.Connection, gate.sessionDidConnected)
	gate.On(bi.Disconnection, gate.sessionDidDisconnected)
	gate.On("publish", gate.sessionWillPublish)
	gate.On("request", gate.sessionWillRequest)
}

func (gate *Gate) sessionDidConnected(sess *Session) {
	id := sess.GetID()
	err := gate.ns.Request2(SubjectID_SessionDidConnect.String(), &Connect{ID: id, Gate: gate.id}, nil, time.Second*5)
	if nil == err {
		idx := len(id) + 1
		sess.Sub, err = gate.ns.Subscribe2(id+".*", idx, func(method string, data []byte) {
			//notice
			// sess.MarshalledEvent(data)
		})
	}
}

func (gate *Gate) sessionDidDisconnected(sess *Session) {
	id := sess.GetID()
	nc, err := nats.Connect(nats.DefaultURL)
	if nil == err {
		msg := &Disconnect{ID: id, Gate: gate.id}
		v, _ := proto.Marshal(msg)
		nc.Publish(SubjectID_SessionDidDisconnect.String(), v)
	}
	sess.Sub.Unsubscribe()
}

func (gate *Gate) sessionWillPublish(sess *Session, msg *Publish) {
	gate.ns.Publish1(msg.Subject, msg.Data)
}

func (gate *Gate) sessionWillRequest(sess *Session, req *Request) *Response {
	data, err := gate.ns.Request1(req.Subject, req.Data, time.Second*5)
	return &Response{Err: err.Error(), Data: data}
}
