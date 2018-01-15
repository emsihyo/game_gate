package main

import (
	"reflect"
	"time"

	"github.com/emsihyo/bi"
	nats "github.com/nats-io/go-nats"
)

//Nats Nats
type Nats struct {
	url         string
	concurrence chan struct{}
	protocol    bi.Protocol
}

//NewNats NewNats
func NewNats(url string, concurrence int, protocol bi.Protocol) *Nats {
	return &Nats{url: url, concurrence: make(chan struct{}, concurrence), protocol: protocol}
}

//Publish1 Publish1
func (n *Nats) Publish1(subject string, v interface{}) error {
	nc, err := n.Get()
	if nil == err {
		var data []byte
		data, err = n.protocol.Marshal(v)
		if nil == err {
			nc.Publish(subject, data)
		}
	}
	n.Put(nc)
	return err
}

//Request1 Request1
func (n *Nats) Request1(subject string, req interface{}, resp interface{}, timeout time.Duration) error {
	nc, err := n.Get()
	if nil == err {
		var data []byte
		data, err = n.protocol.Marshal(req)
		if nil == err {
			var m *nats.Msg
			m, err = nc.Request(subject, data, timeout)
			if nil == err {
				if nil != resp {
					v := reflect.New(reflect.TypeOf(resp).Elem()).Interface()
					err = n.protocol.Unmarshal(m.Data, v)
				}
			}
		}
	}
	n.Put(nc)
	return err
}

//Subscribe1 Subscribe1
func (n *Nats) Subscribe1(id string, method string, f func(v interface{})) (sub *nats.Subscription, err error) {
	var nc *nats.Conn
	nc, err = n.Get()
	if nil == err {
		sub, err = nc.Subscribe(id+"."+method, func(m *nats.Msg) {
			v := reflect.New(reflect.TypeOf(f).In(1).Elem()).Interface()
			f(v)
		})
	}
	n.Put(nc)
	return sub, err
}

//Subscribe2 Subscribe2
func (n *Nats) Subscribe2(id string, method string, f func(v interface{}) interface{}) (sub *nats.Subscription, err error) {
	var nc *nats.Conn
	nc, err = n.Get()
	if nil == err {
		sub, err = nc.Subscribe(id+"."+method, func(m *nats.Msg) {
			v := reflect.New(reflect.TypeOf(f).In(1).Elem()).Interface()
			data, _ := n.protocol.Marshal(f(v))
			if 0 < len(m.Reply) {
				nc.Publish(m.Reply, data)
			}
		})
	}
	n.Put(nc)
	return sub, err
}

//Get Get
func (n *Nats) Get() (*nats.Conn, error) {
	n.concurrence <- struct{}{}
	return nats.Connect(n.url)
}

//Put Put
func (n *Nats) Put(conn *nats.Conn) {
	if nil != conn {
		conn.Close()
	}
	<-n.concurrence
}
