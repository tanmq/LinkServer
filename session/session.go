package session

import (
	"net"

	"linkServer/packet"
)

//Session store user context.
type Session struct {
	UID    uint64
	Device packet.DeviceType

	Conn   *net.TCPConn // tcp connection
	Closed bool
}

//Close tcp connection
func (ss *Session) Close() {
	ss.Conn.Close()
	ss.Closed = true
}
