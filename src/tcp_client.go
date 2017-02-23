package main

import (
    "net"
    "jathy/src/protocol"
    "fmt"
    test "jathy/src/test"
    "time"
    "code.google.com/p/goprotobuf/proto"
)

func main() {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:9909")
    conn, err := net.DialTCP("tcp4", nil, tcpAddr)  

    if err != nil {
        fmt.Println("DialTCP error:", err)
        return     
    }
    defer conn.Close()
    ltest := &test.LoginTest{
        Name: proto.String("one"),      
        Passwd: proto.String("password"),
    } 

    data, err := proto.Marshal(ltest)
    if err != nil {
        fmt.Println("proto.Marshal err:", err) 
        return
    }

    fmt.Printf("data is %v\n", data)
    
    pkt := protocol.NewPacPacket(data, 1)
    
    pdata := pkt.Serialize()
    fmt.Println("pdata", pdata) 
    plen, err := conn.Write(pdata)
    if err != nil {
        fmt.Println("conn.Write err:", err)     
        return
    }
    if plen != len(pdata) {
        fmt.Println("conn.Write len err ", plen, len(pdata))     
        return
    }
    fmt.Printf("write ok\n")
    pacProto := &protocol.PacProtocol{}
    conn.SetDeadline(time.Now().Add(0 * time.Second))
    rpkg, err := pacProto.ReadPacket(conn)
    if err != nil {
        fmt.Println("pacProto.ReadPacket err 1 ", err)    
    } else {
        fmt.Printf("read new pkg %v\n", rpkg)
        return
    }
    time.Sleep(2*time.Second) 
    fmt.Println("wake")
    conn.SetDeadline(time.Now().Add(1 * time.Second))
    rpkg, err = pacProto.ReadPacket(conn)
    if err != nil {
        fmt.Println("pacProto.ReadPacket err 2 ", err)    

        time.Sleep(3*time.Second) 
        conn.SetDeadline(time.Now().Add(1 * time.Second))
        rpkg, err = pacProto.ReadPacket(conn)
        if err != nil {
            fmt.Println("again read error", err)     
        }
    }
    fmt.Printf("read new pkg %v\n", rpkg)
    return
}
