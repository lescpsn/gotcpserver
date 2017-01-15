package gotcpserver

import (
	"net"
)

type Packet interface {
	PacketToByte() []byte
}

type Protocol interface {
	ReadPacket(conn net.Conn) (Packet, error)
}

type MyPacket struct {
	data []byte
}

type MyProtocol struct {
}

func (mp *MyPacket) PacketToByte() []byte {
	return mp.data
}

func (mp *MyProtocol) ReadPacket(conn net.Conn) (Packet, error) {
	rbuf := make([]byte, MaxByteLen)
	rlen, _ := conn.Read(rbuf)
	return &MyPacket{
		data: rbuf[:rlen],
	}, nil
}
