package gotcpserver

import (
	"log"
	"net"
	"sync"
	"time"
)

// NewServerConnection : 每新来一个连接请求实例化一个ServerConnection
func NewServerConnection(srv *Server, c net.Conn) *ServerConnection {
	sc := &ServerConnection{
		srv:               srv,
		conn:              c,
		packetReceiveChan: make(chan interface{}, ReceiveChanSize),
		packetSendChan:    make(chan interface{}, SendChanSize),
		once:              &sync.Once{},
		connCloseChan:     make(chan struct{}),
		remoteAddr:        c.RemoteAddr().String(),
	}
	srv.connMap.Put(sc.remoteAddr, sc)
	return sc
}

// Do : 每个ServerConnection需要处理的事
func (sc *ServerConnection) Do() {
	go sc.readLoop()
	go sc.writeLoop()
	go sc.handleLoop()
}

func (sc *ServerConnection) readLoop() {
	defer func() {
		recover()
		sc.Close()
	}()
	for {
		select {
		case <-sc.connCloseChan:
			return
		default:
		}

		p, _ := mp.ReadPacket(sc.conn)
		sc.packetReceiveChan <- p
		log.Println("[INFO] Read Data From Client:", p, sc.remoteAddr)
	}
}

func (sc *ServerConnection) writeLoop() {
	defer func() {
		recover()
		sc.Close()
	}()

	for {
		select {
		case <-sc.connCloseChan:
			return
		case p := <-sc.packetSendChan:

			bb, ok := p.(*MyPacket)
			if ok {
				log.Println("[INFO] Send Data To Client:", sc.conn.RemoteAddr().String(), bb, bb.data)
			}

			if _, err := sc.conn.Write(bb.data); err != nil {
				return
			}
		}
	}
}

func (sc *ServerConnection) handleLoop() {
	defer func() {
		recover()
		sc.Close()
	}()
	for {
		select {
		case <-sc.connCloseChan:
			return
		case p := <-sc.packetReceiveChan:

			aa, ok := p.(*MyPacket)
			if ok {
				log.Println("[INFO] Wait Handle*************880k:", aa)
			}
			log.Println("[INFO] Wait Handle:", sc.remoteAddr)

			sc.AsyncWritePacket(aa, 2)
		}
	}
}

// AsyncWritePacket : async write packete to client
func (sc *ServerConnection) AsyncWritePacket(p Packet, timeout time.Duration) (err error) {
	defer func() {
		recover()
	}()

	if timeout == 0 {
		select {
		case sc.packetSendChan <- p:
			return nil
		default:
			return ErrWriteBlock
		}
	} else {
		select {
		case sc.packetSendChan <- p:
			return nil
		case <-time.After(timeout):
			return ErrWriteBlock
		}
	}
}

// Close : 关闭ServerConnection连接
func (sc *ServerConnection) Close() {
	sc.once.Do(func() {
		sc.conn.Close()
		close(sc.connCloseChan)
		close(sc.packetSendChan)
		close(sc.packetReceiveChan)
		sc.srv.connMap.Delete(sc.conn.RemoteAddr().String())
		//log.Println("[INFO] ****", sc.srv.connMap)
	})
}

// NewServerConnectionMap : ServerConnectionMap实例化
func NewServerConnectionMap() *ServerConnectionMap {
	return &ServerConnectionMap{
		scm: make(map[string]*ServerConnection),
	}
}

// Put : 新增或者修改ServerConnectionMap中的一个键值对k,v
func (scm *ServerConnectionMap) Put(k string, v *ServerConnection) {
	scm.Lock()
	scm.scm[k] = v
	scm.Unlock()
}

// Delete : 删除ServerConnectionMap中的一个键值对k,v
func (scm *ServerConnectionMap) Delete(k string) {
	scm.Lock()
	delete(scm.scm, k)
	scm.Unlock()
}

// Get : 获取ServerConnectionMap中的一个键值对
func (scm *ServerConnectionMap) Get(k string) (*ServerConnection, bool) {
	scm.Lock()
	sc, ok := scm.scm[k]
	scm.Unlock()
	return sc, ok
}
