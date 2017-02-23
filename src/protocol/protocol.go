package protocol

import (
    "encoding/binary"
    "net"
    "errors"
    "io"
    "fmt"
    "time"
)

const (
    PAC_HEAD_LEN = 8
    PAC_LIMIT_LEN = 10240
)



type PacPacket struct {
    Len     uint32  //长度
    Msgid   uint16  //协议id
    Protp   uint16  //协议类型

    Data    []byte
} 

func (p *PacPacket) Serialize() []byte {
    sb := make([]byte, PAC_HEAD_LEN+len(p.Data))     
    binary.BigEndian.PutUint32(sb[0:4], uint32(len(p.Data))) 
    binary.BigEndian.PutUint16(sb[4:6], p.Msgid) 
    binary.BigEndian.PutUint16(sb[6:8], p.Protp) 
    copy(sb[8:], p.Data)
    return sb
}

func NewPacPacket(buff []byte, id uint16) *PacPacket {
    l := len(buff)
    p := &PacPacket{
        Len: uint32(l),
        Msgid: id,

        Data: buff,
    }
    return p
}

type PacProtocol struct {
}

func (pt *PacProtocol) ReadPacket(conn *net.TCPConn) (*PacPacket, error) {    
    var (
        head []byte = make([]byte, PAC_HEAD_LEN)  
        length uint32 
        mid    uint16
        ptp    uint16
    )
    //_ = time.Now()
    conn.SetDeadline(time.Now().Add(time.Second))
    if _, err := io.ReadFull(conn, head); err != nil {
        return nil, err      
    }

    if length = binary.BigEndian.Uint32(head[0:4]); length > PAC_LIMIT_LEN {
        return nil, errors.New("size of packet is larger than the limit")
    }

    mid = binary.BigEndian.Uint16(head[4:6])
    ptp = binary.BigEndian.Uint16(head[6:8])
    if 0 != ptp {
        fmt.Println("ptp error: ", ptp)    
    }
    
    data := make([]byte, length)
    if _, err := io.ReadFull(conn, data); err != nil {
        return nil, err      
    }

    return NewPacPacket(data, mid), nil
}


