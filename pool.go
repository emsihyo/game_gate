package main

import "sync"

//Pool Pool
var Pool = newPool()

type pool struct {
	Notice *poolNotice
	// Response *
}

func newPool() *pool {
	return &pool{Notice: newPoolNotice()}
}

type poolNotice struct {
	sync.Pool
}

func newPoolNotice() *poolNotice {
	return &poolNotice{Pool: sync.Pool{New: func() interface{} { return &Notice{} }}}
}

func (p *poolNotice) Get() *Notice {
	return p.Pool.Get().(*Notice)
}
func (p *poolNotice) Put(v *Notice) {
	// v.Subject = ""
	// v.Data = nil
	p.Pool.Put(v)
}

type pollResponse struct {
	sync.Pool
}

func newPoolResponse() *pollResponse {
	return &pollResponse{Pool: sync.Pool{New: func() interface{} { return &Response{} }}}
}

func (p *pollResponse) Get() *Response {
	return p.Pool.Get().(*Response)
}
func (p *pollResponse) Put(v *Response) {
	// v.Err = ""
	// v.Data = nil
	p.Pool.Put(v)
}
