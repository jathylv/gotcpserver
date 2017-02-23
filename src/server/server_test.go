package svr

import (
    "net"
    "sync/atomic"
    "testing"
    "fmt"
    "time"
    "jathy/src/protocol"
    "jathy/src/test"
    "code.google.com/p/goprotobuf/proto"
)

var conn_id uint32 = 0

type TConn struct {
    conn *net.TCPConn  
    closeChan chan int
    recvChan chan interface{}
    writeChan chan interface{}
    closeFlag int32
    addr string
    id  uint32
}



func (c *TConn) GetRawConn() *net.TCPConn {
    return c.conn     
}

func (c *TConn) CloseFlag() bool {
    return atomic.LoadInt32(&c.closeFlag) == 1
}

func (c *TConn) SetCloseFlag() bool {
    b := atomic.CompareAndSwapInt32(&c.closeFlag, 0, 1)  
    if b {
        close(c.closeChan) 
    }
    return b
}

func (c *TConn) GetCloseChan() chan int {
    return c.closeChan    
}

func (c *TConn) GetRecvChan() chan interface{} {
    return c.recvChan    
}

func (c *TConn) GetWriteChan() chan interface{} {
    return c.writeChan    
}


type TSvr struct {
    cm map[uint32]interface{}
}

func (t *TSvr) OnConnect(conn TcpConnItf) bool {
    if nc, ok := conn.(*TConn); ok {
        fmt.Printf("OnConnect new connect id %d from: %s\n", nc.id, nc.addr)     
        t.cm[nc.id] = nc
        return true
    } else {
        fmt.Println("OnConnect not ok")      
        return false
    }
}

func (t *TSvr) OnMessage(conn TcpConnItf, pkg interface{}) bool {
    fmt.Printf("OnMessage %s\n", conn.GetRawConn().RemoteAddr().String())
    if nc, ok := conn.(*TConn); ok {
        fmt.Printf("OnMessage new msg from: %s data %v\n", nc.addr, pkg)     
        if pp, ok := pkg.(*protocol.PacPacket); ok {
            if 1 == pp.Msgid {
                ltest := &test.LoginTest{}       
                if err := proto.Unmarshal(pp.Data, ltest); err != nil {
                    fmt.Printf("proto.Unmarshal err %v\n", err)     
                } else {
                    fmt.Printf("Name %s, Passwd %s\n", ltest.GetName(), ltest.GetPasswd()) 
                } 
            }  
        } else {
            fmt.Printf("OnMessage type error\n")  
        }
        nc.writeChan <- pkg
        return true
    } else {
        fmt.Println("OnMessage not ok")      
        return false
    }
}

func (t *TSvr) NewConn(conn *net.TCPConn, len int) TcpConnItf {
    tc := &TConn{
        conn: conn,
        closeFlag: 0,
        closeChan: make(chan int),
        recvChan: make(chan interface{}, len),
        writeChan: make(chan interface{}, len),
        addr: conn.RemoteAddr().String(),
    } 
    tc.id = atomic.AddUint32(&conn_id, 1) 
    if true == atomic.CompareAndSwapUint32(&conn_id, 0, 1) {
        tc.id = 1     
    }

    return tc
}

func (t *TSvr) OnClose(conn TcpConnItf) {
    if nc, ok := conn.(*TConn); ok {
        fmt.Printf("OnClose connect id %d from: %s\n", nc.id, nc.addr)     
        if _, ok := t.cm[nc.id]; ok {
            fmt.Println("OnClose ok:", nc.addr)
            delete(t.cm, nc.id)  
        } else {
            fmt.Println("OnClose delete error")
        }
    } else {
        fmt.Println("OnClose not ok")      
    }
}

func TestSvr(t *testing.T) {
    ts := &TSvr{
        cm: make(map[uint32]interface{}),     
    }
    s := NewTcpSvr("127.0.0.1:9909", ts, 100) 
    
    go s.Start()
    fmt.Println("main start")
    <- time.After(60*time.Second) 
    fmt.Println("main stop")
    s.Stop()
}

