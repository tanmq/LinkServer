package util

import (
  "net"
  "errors"

  "linkServer/packet"
)


const (
  MAX_DATA_Length = 1 << 23  //data should less than 1MB
)

//数据格式开始四个byte表示data长度，后面跟随data信息。这里会接收一个完整的data包才会返回
//或者读取出错也会返回，出错后会断开该链接
func FixedLengthRead(conn *net.TCPConn) (data []byte, err error) {
  //start to read data length
  bufLen, err := readNBytes(conn, 4)
  if err != nil {
    return
  }
  //start to read real data
  datalen := packet.BytesToUint32(bufLen)
  if datalen > MAX_DATA_Length {
    err = errors.New("Data are too large.")
    return
  }

  data, err = readNBytes(conn, datalen)

  return
}


func readNBytes(conn *net.TCPConn, num uint32) (data []byte, err error) {
   data = make([]byte, 0, num)
   var nread  uint32    = 0
   var unread uint32    = num

   for nread < num {
     b := make([]byte, unread)
     c, e := conn.Read(b)
     if e != nil {
       err = e
       return
     }
     data = append(data, b[:c]...)
     nread  += uint32(c)
     unread -= uint32(c)
   }

   if uint32(len(data)) != num {
     err = errors.New("Read error. Cannot reach enough bytes.")
   }

   return
}



func FixedLengthWrite(conn *net.TCPConn, data []byte) error {
  datalen := uint32(len(data))
  _, err := conn.Write(packet.Uint32ToBytes(datalen))
  if err != nil {
    return err
  }

  _, err = conn.Write(data)
  return err
}
