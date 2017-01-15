package main

import (
	"iota/dtusimulator"
	"log"
	"net"
	"strconv"
	"sync"
)

var (
	maxConnNum = 10
	address    = "127.0.0.1:9012"
)

type MyDtu struct {
}

func (md *MyDtu) OnReceive(ibuf []byte) (obuf []byte) {
	ilen := len(ibuf)
	if ilen > 0 {
		buf := make([]byte, ilen)
		buf[0] = 'a'
		return buf[:1]
	} else {
		return nil
	}

}

func StartConn(maxConnNum int) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	nc := make([]net.Conn, maxConnNum)
	var err error
	for i := 0; i < maxConnNum; i++ {
		nc[i], err = net.Dial("tcp", address)
		if err != nil {
			log.Println("[ERRO] Net Dail Error:", err)
			return
		}

		log.Println("[INFO] Connect Info:", nc[i].LocalAddr())
	}

	log.Println("[INFO] Connect Info:", maxConnNum, len(nc), nc)
	for i := 0; i < maxConnNum; i++ {
		md := &MyDtu{}
		id := strconv.Itoa(10000000 + i)
		phoneNumber := strconv.Itoa(18600000000 + i)
		ip := "127.0.0.1"
		dtu := dtusimulator.NewDtu(id, phoneNumber, ip, md)
		wg.Add(1)
		go dtu.Start(nc[i])
	}
	wg.Wait()
}

func main() {
	StartConn(maxConnNum)
}
