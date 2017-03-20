package server


import (
	"net"
	"strconv"
	"os"
	"time"
	"errors"


	"linkServer/logger"
	"linkServer/session"
	"linkServer/packet"
	"linkServer/common/util"
)



const (
	MAX_EXPIRE_SEC = 30
)

type TcpServer struct {
	sessions session.SessionMap
}


func NewTcpServer() *TcpServer {
	server := new(TcpServer)
	server.sessions = session.NewSessionMap()
	return server
}

func (s *TcpServer)report() {
	ticker := time.NewTicker(time.Millisecond * 10)
	go func() {
    for _ = range ticker.C {
        logger.Info("ticked at ", time.Now(), ". The size of sessions is ",  s.sessions.SessionSize())
    }
	}()
}



func (s *TcpServer)ListenAndServe(port int)  error {
  tcpAddr, err := net.ResolveTCPAddr("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		panic("init tcp addr error")
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Error("Error to bind on port : %d", port)
		os.Exit(-1)
	}

	defer l.Close()

	s.report()

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			logger.Error("Error Accepted. ", err)
			continue
		}

		go s.handleLink(conn)
	}
}



func (s *TcpServer)handleLink(tcpConn *net.TCPConn) {

	// 1s内需要收到登录包
	tcpConn.SetDeadline(time.Now().Add(time.Second))
	logger.Info("Accept new Link!")

	//build session
	ss, err := s.buildSession(tcpConn)
	if (err != nil) {
		logger.Error("Error to build session. so connection is closing.", err)
		tcpConn.Close()
		return
	}

	//read data
	//设置读时间超时30s，正常情况下30s内肯定会有心跳包的到达，如果有数据传输则代替心跳包
	tcpConn.SetDeadline(time.Now().Add(time.Second * MAX_EXPIRE_SEC))
	for {
		data, err := util.FixedLengthRead(tcpConn)
		if (err != nil) {
			logger.Error("something wrong to connection.", err)
			break
		}
		pack, err := packet.NewPacket(data)
		if err != nil {
			logger.Warn("broken packet")
			continue
		}
		s.dispatch(pack)
	}

	//clear connection and session
	if !ss.Closed {
		s.sessions.Del(ss.Uid, ss.Device)
	}
}

//处理数据包
func (s *TcpServer)dispatch(pack packet.Packet) {
	 switch pack.Protocol {
	 case packet.HEART_BEAT:
		 return //ignore


	 }
}



//build user session and join into sessionmap
func (s *TcpServer)buildSession(tcpConn *net.TCPConn) (ss *session.Session, err error) {
	 data, err := util.FixedLengthRead(tcpConn)
	 if err != nil {
		 return
	 }

	 pack, err 		:= packet.NewPacket(data)
	 if err != nil {
		 	logger.Error("Cannot build packet.")
		 	return
	 }
	 loginInfo, err := pack.ParseLoginInfo()
	 if err != nil {
		 return
	 }

	 //登录验证
	 if (!util.Auth(&loginInfo)) {
		 err = errors.New("auth failed.")
		 return
	 }


	 //记录session信息
	 ss 				= new(session.Session)
	 ss.Uid 		= loginInfo.Uid
	 ss.Closed 	= false
	 switch packet.DeviceType(loginInfo.Device) {
	 	case packet.APP:
		 	ss.Device = packet.APP
		case packet.WEB:
			ss.Device = packet.WEB
		default:
			ss.Device = packet.UNKNOW
	 }
	 ss.Conn = tcpConn

	 // add to session list
	 s.sessions.Put(ss.Uid, ss.Device, ss)

	 return
}
