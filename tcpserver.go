package gotcpserver

import (
	"log"
	"net"
	"os"
)

var (
	mp = &MyProtocol{} // 等待完善
)

// NewServer : Server实例化
func NewServer(addr string) *Server {
	return &Server{
		laddr:    addr,
		connMap:  NewServerConnectionMap(),
		exitChan: make(chan struct{}),
	}
}

// Start : 启动服务器
func (srv *Server) Start() {
	l, err := net.Listen("tcp", srv.laddr)
	if err != nil {
		log.Println("[ERRO] Server Listening Error:", err)
		os.Exit(-1)
	}

	defer func() {
		l.Close()
	}()

	log.Println("[INOF] Server Is Listening At:", srv.laddr)
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("[ERRO] Accept Error:", err)
			continue
		}
		log.Println("[INOF] Remote Host Login:", c.RemoteAddr().String())

		go NewServerConnection(srv, c).Do()
	}
}

// Stop : 停止服务器
func (srv *Server) Stop() {
	log.Println("[INFO] Stopping Server...")
	close(srv.exitChan)
}
