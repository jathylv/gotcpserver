package svr

import (
    "net"
)

type TcpConnItf interface {
    GetRawConn() *net.TCPConn 
    CloseFlag() bool 
    SetCloseFlag() bool
    GetCloseChan() chan int
    GetRecvChan() chan interface{}
    GetWriteChan() chan interface{}
} 
