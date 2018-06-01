package util

import (
	"errors"
	"net"

	"linkServer/packet"
)

const (
	//MaxDataLength means the longest data length.
	MaxDataLength = 1 << 23 //data should less than 1MB
)

//FixedLengthRead read fixedLength data.When not reach the length, it will block.
func FixedLengthRead(conn net.Conn) (data []byte, err error) {
	//start to read data length
	bufLen, err := readNBytes(conn, 4)
	if err != nil {
		return
	}
	//start to read real data
	datalen := packet.BytesToUint32(bufLen)
	if datalen > MaxDataLength {
		err = errors.New("data to large")
		return
	}

	data, err = readNBytes(conn, datalen)

	return
}

func readNBytes(conn net.Conn, num uint32) (data []byte, err error) {
	data = make([]byte, 0, num)
	var nread uint32 = 0
	var unread = num

	for nread < num {
		b := make([]byte, unread)
		c, e := conn.Read(b)
		if e != nil {
			err = e
			return
		}
		data = append(data, b[:c]...)
		nread += uint32(c)
		unread -= uint32(c)
	}

	if uint32(len(data)) != num {
		err = errors.New("cannot reach enough bytes")
	}

	return
}

//FixedLengthWrite write fixedLength data to conn.
func FixedLengthWrite(conn *net.TCPConn, data []byte) error {
	datalen := uint32(len(data))
	bufData := packet.Uint32ToBytes(datalen)
	bufData = append(bufData, data...)
	_, err := conn.Write(bufData)
	return err
}
