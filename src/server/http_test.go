package svr

import (
    "jathy/src/protocol"
    "jathy/src/test"
    "code.google.com/p/goprotobuf/proto"
    "net/http"
    "fmt"
    "testing"
    "time"
)

func HttpHandleFunc(w http.ResponseWriter, mid uint16, data []byte) {
    fmt.Printf("new http msg %v %v\n", mid, data)
    if 2 == mid {
        msg := &test.HttpMsg{}        
        if err := proto.Unmarshal(data, msg); err != nil {
            fmt.Printf("proto.Unmarshal err %v\n", err)     
        } else {
            fmt.Printf("method %s, userid %v, data %v\n", msg.GetMethod(), msg.GetUserid(), msg.GetData()) 
        } 
        p := protocol.NewPacPacket(data, mid)
        fmt.Fprintf(w, "%v", p.Serialize())
    }     
}

func TestHttpSvr(t *testing.T) {
    hs := NewMyHttpSvr() 
    fmt.Println("hs StartFunc")
    go hs.StartFunc(":9710", HttpHandleFunc)
    <- time.After(20 * time.Second)
}

