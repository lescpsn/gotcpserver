package gotcpserver

import (
	"errors"
)

const (
	MaxByteLen      = 4096
	SendChanSize    = 1024
	ReceiveChanSize = 1024
)

var (
	ErrWriteBlock = errors.New("Write Packet Is Blocked Time Out")
	ErrReadBloc   = errors.New("Read Packet Is Blocked Time Out")
)
