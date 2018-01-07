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
func (n *Nats) Publish1(subject string, data []byte) error {
	nc, err := n.Get()
	if nil == err {
		nc.Publish(subject, data)
	}
	n.Put(nc)
	return err
}

//Publish2 Publish2
func (n *Nats) Publish2(subject string, v interface{}) error {
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
func (n *Nats) Request1(subject string, data []byte, timeout time.Duration) ([]byte, error) {
	nc, err := n.Get()
	var resp []byte
	if nil == err {
		var m *nats.Msg
		m, err = nc.Request(subject, data, timeout)
		if nil == err {
			resp = m.Data
		}
	}
	n.Put(nc)
	return resp, err
}

//Request2 Request2
func (n *Nats) Request2(subject string, req interface{}, resp interface{}, timeout time.Duration) error {
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
func (n *Nats) Subscribe1(subject string, f func(v interface{})) (*nats.Subscription, error) {
	var sub *nats.Subscription
	nc, err := n.Get()
	if nil == err {
		nc.Subscribe(subject, func(m *nats.Msg) {
			v := reflect.New(reflect.TypeOf(f).In(0).Elem()).Interface()
			f(v)
		})
	}
	n.Put(nc)
	return sub, err
}

//Subscribe2 Subscribe2
func (n *Nats) Subscribe2(subject string, offset int, f func(method string, data []byte)) (*nats.Subscription, error) {
	var sub *nats.Subscription
	nc, err := n.Get()
	if nil == err {
		nc.Subscribe(subject, func(m *nats.Msg) {
			f(subject[offset:], m.Data)
		})
	}
	n.Put(nc)
	return sub, err
}

//Subscribe3 Subscribe3
func (n *Nats) Subscribe3(subject string, offset int, f func(method string, v interface{})) (*nats.Subscription, error) {
	var sub *nats.Subscription
	nc, err := n.Get()
	if nil == err {
		nc.Subscribe(subject, func(m *nats.Msg) {
			v := reflect.New(reflect.TypeOf(f).In(1).Elem()).Interface()
			f(subject[offset:], v)
		})
	}
	n.Put(nc)
	return sub, err
}

//Subscribe4 Subscribe4
func (n *Nats) Subscribe4(subject string, offset int, f func(method string, v interface{}) interface{}) (*nats.Subscription, error) {
	var sub *nats.Subscription
	nc, err := n.Get()
	if nil == err {
		nc.Subscribe(subject, func(m *nats.Msg) {
			v := reflect.New(reflect.TypeOf(f).In(1).Elem()).Interface()
			data, _ := n.protocol.Marshal(f(subject[offset:], v))
			nc.Publish(m.Reply, data)
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
