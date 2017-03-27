package packet

import (
	"encoding/binary"
	"encoding/json"
	"errors"
)

//DeviceType is user login device. app or web.
type DeviceType uint8

//ProtocolType is used to deliver packet.
type ProtocolType uint32

const (
	//APP means ios or android
	APP DeviceType = 1
	//WEB means web login
	WEB DeviceType = 2
	//UNKNOW means can not verify.
	UNKNOW DeviceType = 8
)

const (
	//HeartBeat protocol used to keep communicate will client.
	HeartBeat ProtocolType = iota
	//Login protocol used to login.
	Login
)

//Packet represent data received from client.
type Packet struct {
	data     []byte
	Protocol ProtocolType
}

//NewPacket is to build byte stream into packet.
func NewPacket(data []byte) (Packet, error) {
	packet := Packet{}
	packet.data = data

	e := packet.parse()
	if e != nil {
		return packet, e
	}

	return packet, nil
}

//ParseLoginInfo is to parse login data from packet
func (p *Packet) ParseLoginInfo() (LoginInfo, error) {
	if ProtocolType(p.Protocol) != Login {
		return LoginInfo{}, errors.New("not an login packet")
	}

	info := LoginInfo{}
	e := json.Unmarshal(p.data[4:], &info)
	if e != nil {
		return info, e
	}

	return info, nil
}

func (p *Packet) parse() error {
	if len(p.data) < 4 {
		return errors.New("unable to parse")
	}

	p.Protocol = ProtocolType(BytesToUint32(p.data[:4]))

	return nil
}

//Uint64ToBytes is to encode uint64 into bytes use BigEndian
func Uint64ToBytes(n uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, n)

	return bytes
}

//BytesToUint64 is to decode bytes into uint64 use BigEndian
func BytesToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

//Uint32ToBytes is to encode uint32 into bytes use BigEndian
func Uint32ToBytes(n uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, n)

	return bytes
}

//BytesToUint32 is to decode bytes into uint32 use BigEndian
func BytesToUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}
