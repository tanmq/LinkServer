package tcpnode

import (
	"time"
	"net"
	"fmt"

	"linkServer/packet"
	"linkServer/common/util"
)


const WAIT_FRESHOLD_TIME = 60


type NodeClient struct {
	receiver 	MessageReceive

	maxmessagelength  uint32
	address						string
	conn	  					net.Conn
	connected					bool
}


func NewNodeClient(address string ,sender MessageSender, receiver MessageReceiver) *NodeClient {
	client := new(NodeClient)

	client.sender						 	= sender
	client.receiver				 		= receiver
	client.address						= address
	client.connected					= false

	client.connect()

	return client
}

func (c *NodeClient)IsValid() {
	return c.connected
}


//reconnect the tcp connection.
func (c *NodeClient)Reset() {
	c.conneced = false
	if c.conn != nil { // close old connection
		c.conn.Close()
		fmt.Println("We are closing the old connection and begin to reconnect for address : ", c.address)
	}
	connect()
}

func (c *NodeClient)connect() {

	waitTime := 1
	for !connected {
		c.conn, err = net.Dial("tcp", c.address)
		if err == nil {
			c.conn.SetWriteDeadline(time.Now().Add(time.Second))
			c.connected = true
			go c.handleReceive()
			break
		}

		if waitTime >= MAX_EXPIRE_SEC {
			waitTime = MAX_EXPIRE_SEC
		} else {
			waitTime = 2 * waitTime
		}

		fmt.Println("Connect failed. We are now going to sleep and then retry. waitTime is ", waitTime)
		time.Sleep(waitTime*time.Second)
	}

	return
}


func (c *NodeClient)Send(data byte[]) error {
	 datalen := len(data)
	 if datalen > c.maxMessageLength {
		 return errors.New("Data is too big.")
	 }

	 lenBuf := packet.Uint32ToBytes(uint32(datalen))
	 data = append(lenBuf, data...)
	 _, err := c.conn.Write(data)
	 if err != nil {
		 go Reset()
	 }

	 return err
}



func (c *NodeClient)handleReceive() {
	for c.connected {
		data, err := ioutil.FixedLengthRead(c.conn)
		if err != nil {
			fmt.Println("Wrong to read data. And we will receive data wait some time.")
		}

		receiver.OnMessage(data)
	}
}
