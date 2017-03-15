package svr

import (
    "net/http"
    "jathy/src/protocol"
    "fmt"
    "log"
)

var thpPro protocol.PacProtocol
var  hf func(http.ResponseWriter, uint16, []byte)

type HttpMySvr struct {
    addr string
}

func NewMyHttpSvr() *HttpMySvr {
    return &HttpMySvr{
        addr: "",
    }    
}

func myHandle(w http.ResponseWriter, r *http.Request) {
    fmt.Println("StartFunc in")
    log.Println("StartFunc in")
    p, err := thpPro.ReadHttpPacket(r)  
    if err != nil {
        fmt.Println(err)          
        fmt.Fprintf(w, "err %v", err)
        return
    }
    hf(w, p.Msgid, p.Data)
    return
}

func (s *HttpMySvr) StartFunc(addr string, hfunc func(http.ResponseWriter, uint16, []byte)) {
    hf = hfunc
    http.HandleFunc("/", myHandle)  
    log.Fatal(http.ListenAndServe(addr, nil))
}

