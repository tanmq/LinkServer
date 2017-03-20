package session

import (
  "net"

  "linkServer/packet"
)


type Session struct {
  Uid             uint64
  Device          packet.DeviceType

  Conn            *net.TCPConn // tcp connection
  Closed          bool
}


//close tcp connection
func (ss *Session)Close() {
  ss.Conn.Close()
  ss.Closed = true
}
