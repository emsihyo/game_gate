package main

import "net"
import "github.com/emsihyo/bi"
import "flag"
import "time"

func main() {
	addrPtr := flag.String("addr", ":12321", "addr")
	servTCP(*addrPtr)
}

func servTCP(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if nil != err {
		panic(err)
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		panic(err)
	}
	protocol := &bi.ProtobufProtocol{}
	gate := NewGate()
	for {
		tcpConn, err := l.AcceptTCP()
		if nil == err {
			go func() {
				conn := bi.Pool.TCPConn.Get(tcpConn)
				sess := NewSession(conn, protocol, time.Second*20)
				gate.Handle(sess)
			}()
		}
	}
}
