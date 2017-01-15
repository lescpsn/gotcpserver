package gotcpserver

import (
	"net"
	"sync"
)

// Server : Server服务器结构
type Server struct {
	laddr    string
	connMap  *ServerConnectionMap
	exitChan chan struct{}
}

// ServerConnection : Server维持的Connection结构
type ServerConnection struct {
	srv               *Server
	conn              net.Conn
	packetSendChan    chan interface{}
	packetReceiveChan chan interface{}
	once              *sync.Once
	connCloseChan     chan struct{}
	remoteAddr        string //remoteip:port
}

// ServerConnectionMap : map[k]v:
type ServerConnectionMap struct {
	sync.RWMutex
	scm map[string]*ServerConnection
}
