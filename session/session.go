package session

import (
  "net"

  "linkServer/packet"
)

const (
  expireLimit = 30 // second
)

type Session struct {
  Uid             uint64
  Device          packet.DeviceType

  Conn            *net.TCPConn // tcp connection
}
