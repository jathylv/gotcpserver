package svr

import (
    "sync"
    "net"
    "log"
    "time"
    "jathy/src/protocol"
    "fmt"
)

var tPro protocol.PacProtocol

type TcpSvrItf interface {
    OnConnect(conn TcpConnItf) bool      
    OnMessage(conn TcpConnItf, pkg interface{}) bool
    NewConn(conn *net.TCPConn, len int) TcpConnItf
    OnClose(conn TcpConnItf)
}

type TcpSvr struct {
    addr    string    
    chanLen int 
    exitChan chan int
    svrItf  TcpSvrItf
    waitGroup *sync.WaitGroup
}

func NewTcpSvr(addr string, itf TcpSvrItf, len int) *TcpSvr {
    return &TcpSvr{
        addr: addr,     
        chanLen: len,
        exitChan: make(chan int),
        svrItf: itf,
        waitGroup: &sync.WaitGroup{},
    }    
}

func (s *TcpSvr) Start() {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", s.addr)
    checkErr(err)
    listener, err := net.ListenTCP("tcp", tcpAddr) 
    checkErr(err)

    s.waitGroup.Add(1) 

    defer func() {
        listener.Close()
        s.waitGroup.Done()  
    }() 
    fmt.Println("tcpsvr start")
    for {
        select {
        case <- s.exitChan:
            return
        default:

        }
        listener.SetDeadline(time.Now().Add(time.Second))
        newConn, err := listener.AcceptTCP()
        if err != nil {
            log.Println(err)
            continue     
        }

        go s.handleNewConn(newConn)
    }
}

func (s *TcpSvr) Close(conn TcpConnItf) {
    if conn != nil && (true == conn.SetCloseFlag()) {
        fmt.Printf("Close %s\n", conn.GetRawConn().RemoteAddr().String())
        conn.GetRawConn().Close()    
        s.svrItf.OnClose(conn)
    }
}

func (s *TcpSvr) Stop() {
    fmt.Println("tcpsvr stop")
    close(s.exitChan)        
    s.waitGroup.Wait()
}

func (s *TcpSvr) handleNewConn(conn *net.TCPConn) {
    s.waitGroup.Add(1) 
    defer s.waitGroup.Done()
    
    newConn := s.svrItf.NewConn(conn, s.chanLen)

    if nil == newConn {
        conn.Close()
        return     
    }
    defer s.Close(newConn)    

    if !s.svrItf.OnConnect(newConn) {
        return     
    }

    // read proc
    go s.handleReadConn(newConn)
    // write proc
    go s.handleWriteConn(newConn) 
   
    for {
        fmt.Printf("handleNewConn tPro.ReadPacket from %s\n", conn.RemoteAddr().String())
        select {
        case <- s.exitChan:
            return
        case <- newConn.GetCloseChan():
            return
        default:
        } 
        
        if newConn.CloseFlag() {
            return     
        }
        //conn.SetDeadline(time.Now().Add(time.Second))
        p, err := tPro.ReadPacket(conn)
        if err != nil {
            fmt.Printf("tPro.ReadPacket err %v %s\n", err, conn.RemoteAddr().String())
            if operr, ok := err.(*net.OpError); ok && operr.Timeout() {
                continue    
            } else {
                return
            }
        }
        
        fmt.Printf("tPro.ReadPacket p %v from %s\n", p, conn.RemoteAddr().String())
        newConn.GetRecvChan() <- p
    } 

}


func (s *TcpSvr) handleReadConn(conn TcpConnItf) {
    s.waitGroup.Add(1)
    defer s.waitGroup.Done()   
    for {
        select {
        case <- s.exitChan:
            return
        case <- conn.GetCloseChan():
            return
        case p := <- conn.GetRecvChan():
            fmt.Printf("handleReadConn %s\n", conn.GetRawConn().RemoteAddr().String())
            s.svrItf.OnMessage(conn, p)

        } 
    }
}

func (s *TcpSvr) handleWriteConn(conn TcpConnItf) {
    s.waitGroup.Add(1)
    defer s.waitGroup.Done()   
    for {
        select {
        case <- s.exitChan:
            return
        case <- conn.GetCloseChan():
            return
        case p := <- conn.GetWriteChan():
            if conn.CloseFlag() {
                return 
            } 
            fmt.Printf("handleWriteConn p %v %s\n", p, conn.GetRawConn().RemoteAddr().String())
            if pp, ok := p.(*protocol.PacPacket); ok {
                rawConn := conn.GetRawConn()
                rawConn.SetWriteDeadline(time.Now().Add(time.Second))
                rawConn.Write(pp.Serialize())
            } else {
                log.Fatal("not PacPacket type")      
                return
            } 
        } 
    }
       
}

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)      
    }     
}

