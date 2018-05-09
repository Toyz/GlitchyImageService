package engine

import "net"

type Server struct {
	tcpAddr *net.TCPAddr
}

func NewServer(addr string) *Server {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	checkError(err)

	return &Server{
		tcpAddr,
	}
}

func (srv *Server) Listen(runner func(*net.TCPListener)) {
	listener, err := net.ListenTCP("tcp", srv.tcpAddr)
	checkError(err)

	runner(listener)
}
