package packet

import (
  "errors"
  "encoding/binary"
  "encoding/json"

)

type DeviceType uint8
type ProtocolType   uint32

const (
  APP DeviceType    = 1
  WEB DeviceType    = 2
  UNKNOW DeviceType = 8
)


const (
  HEART_BEAT ProtocolType = iota // zero means heart beat packet
  LOGIN                      // login packet

)

//一个packet至少要有4个字节，其中
//  4byte (protocol)
type Packet struct {
  data      []byte
  Protocol  ProtocolType
}

func NewPacket(data []byte)(Packet, error) {
  packet := Packet{}
  packet.data = data

  e := packet.parse()
  if e != nil {
    return packet, e
  }

  return packet, nil
}

func (p *Packet)ParseLoginInfo() (LoginInfo, error) {
  if ProtocolType(p.Protocol) != LOGIN {
    errors.New("Error. Not an login packet.")
  }

  info := LoginInfo{}
  e := json.Unmarshal(p.data[4:], &info)
  if e != nil {
    return info, e
  }

  return info, nil
}


func (p* Packet)parse() error{
   if len(p.data) < 4 {
     return errors.New("Wrong packet. Unable to parse.")
   }

   p.Protocol = ProtocolType(BytesToUint32(p.data[:4]))

   return nil
}


//网络采用的是大端规则，所以这里使用大端规则来控制字节序转换
func Uint64ToBytes(n uint64) []byte {
  bytes := make([]byte, 8)
  binary.BigEndian.PutUint64(bytes, n)

  return bytes
}

func BytesToUint64(bytes []byte) uint64 {
  return binary.BigEndian.Uint64(bytes)
}

func Uint32ToBytes(n uint32) []byte {
  bytes := make([]byte, 4)
  binary.BigEndian.PutUint32(bytes, n)

  return bytes
}

func BytesToUint32(bytes []byte) uint32 {
  return binary.BigEndian.Uint32(bytes)
}
