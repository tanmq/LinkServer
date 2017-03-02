package server


import (
	"net"
	"strconv"
	"os"


	"linkServer/logger"
)

type TcpServer struct {
}


func NewTcpServer() *TcpServer {
	return new(TcpServer)	
}



func (s *TcpServer)ListenAndServe(port int)  error {
	
	l, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		logger.Error("Error to bind on port : %d", port)
		os.Exit(-1)
	}
	
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error Accepted. ", err)
			continue
		}

		go s.handleLink(conn)
	}
}



func (s *TcpServer)handleLink(conn net.Conn) {

	logger.Info("Accept new Link!")
	conn.Close()
}


