package main

import (
    "jathy/src/protocol"
    "fmt"
    test "jathy/src/test"
    "code.google.com/p/goprotobuf/proto"
    "bytes"
    "net/http"
    "os"
    "io"
)

func main() {
    msg := &test.HttpMsg{
        Method: proto.String("myfunc"),
        Userid: proto.Uint32(9527),
        Data:   proto.String("as your like"), 
    } 
    data, err := proto.Marshal(msg)
    if err != nil {
        fmt.Println("proto.Marshal err:", err) 
        return
    }

    fmt.Printf("data is %v\n", data)
    
    pkt := protocol.NewPacPacket(data, 2)
    req, err := http.NewRequest("POST", "http://127.0.0.1:9710/", bytes.NewReader(pkt.Serialize()))
    client := &http.Client{}
    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil {
        fmt.Printf("do err %v\n", err)     
    } else {
        fmt.Printf("do ok \n")     
        io.Copy(os.Stdout, resp.Body)
        fmt.Printf("\n")
    }
    return
}

