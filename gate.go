package main

import (
	"time"

	"github.com/emsihyo/bi"
	nats "github.com/nats-io/go-nats"
)

//Gate Gate
type Gate struct {
	bi.BI
	id     string
	ns     *Nats
	sesses map[string]*Session
}

//NewGate NewGate
func NewGate() *Gate {
	gate := &Gate{ns: NewNats(nats.DefaultURL, 100, &bi.ProtobufProtocol{})}
	return gate
}

func (gate *Gate) init() {
	id := gate.id
	var err error
	_, err = gate.ns.Subscribe1(id, A2G_Kick.String(), gate.reqKick)
	if nil != err {
		panic(err)
	}
	_, err = gate.ns.Subscribe1(id, A2G_Event.String(), gate.reqEvent)
	if nil != err {
		panic(err)
	}
	_, err = gate.ns.Subscribe2(id, A2G_Request.String(), gate.reqRequest)
	if nil != err {
		panic(err)
	}
	_, err = gate.ns.Subscribe2(id, A2G_JoinRoom.String(), gate.reqJoinRoom)
	if nil != err {
		panic(err)
	}
	_, err = gate.ns.Subscribe2(id, A2G_LeaveRoom.String(), gate.reqLeaveRoom)
	if nil != err {
		panic(err)
	}
	gate.On(bi.Connection, gate.sessionDidConnected)
	gate.On(bi.Disconnection, gate.sessionDidDisconnected)

}

func (gate *Gate) sessionDidConnected(sess *Session) {
	id := sess.GetID()
	err := gate.ns.Request1(G2A_SessionDidConnect.String(), &ReqConnect{ID: id, Gate: gate.id}, &RespConnect{}, time.Second*5)
	if nil != err {

	}
}

func (gate *Gate) sessionDidDisconnected(sess *Session) {
	id := sess.GetID()
	err := gate.ns.Request1(G2A_SessionDidConnect.String(), &ReqDisconnect{ID: id, Gate: gate.id}, &RespDisconnect{}, time.Second*5)
	if nil != err {

	}
}

func (gate *Gate) reqEvent(req interface{}) {
	reqEvent := req.(*ReqEvent)
	sess, ok := gate.sesses[reqEvent.ID]
	if ok {
		sess.Event("event", &reqEvent.Data, nil)
	}
}

func (gate *Gate) reqRequest(req interface{}) interface{} {
	resp := RespRequest{}
	reqEvent := req.(*ReqEvent)
	sess, ok := gate.sesses[reqEvent.ID]
	if ok {
		respBytes := []byte{}
		err := sess.Event("event", &reqEvent.Data, &respBytes)
		resp.Data = respBytes
		if nil != err {
			resp.Code = 1
			resp.Desc = err.Error()
		} else {
			resp.Code = 0
		}
	}
	return &resp
}

func (gate *Gate) reqKick(req interface{}) {

}

func (gate *Gate) reqJoinRoom(req interface{}) interface{} {
	return nil
}

func (gate *Gate) reqLeaveRoom(req interface{}) interface{} {
	return nil
}
