package tcpnode

import (
	"time"
	"net"
	"fmt"
	"sync"

	"linkServer/packet"
	"linkServer/common/util"
)


const WAIT_FRESHOLD_TIME_SEC = 600


type NodeClient struct {
	receiver 	MessageReceiver

	maxmessagelength  uint32
	address						string
	conn	  					net.Conn

	mu								sync.Mutex
	connected					bool 			// protected by mu
	canceled					bool  		// protected by mu
}


func NewNodeClient(address string ,receiver MessageReceiver) *NodeClient {
	client := new(NodeClient)

	client.receiver				 		= receiver
	client.address						= address
	client.connected					= false
	client.canceled						= false

	client.connect()

	return client
}

func (c *NodeClient)IsValid()  bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

func (c *NodeClient)valid() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connected = true
}

func (c *NodeClient)unValid() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.conneced = false
}

func (c *NodeClient)IsCancel() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.canceled
}

func (c *NodeClient)Cancel() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.canceled = true
	c.conneced = false
	c.conn.Close()
}


//reconnect the tcp connection.
func (c *NodeClient)Reset() {
	c.unValid()
	if c.conn != nil { // close old connection
		c.conn.Close()
		fmt.Println("We are closing the old connection and begin to reconnect for address : ", c.address)
	}
	connect()
}

/**
* This method will block util the connection is built or the client is canceled.
*/
func (c *NodeClient)connect() {

	waitTime := 1
	for !c.IsValid() && !c.IsCancel() {
		c.conn, err = net.Dial("tcp", c.address)
		if err == nil {
			c.conn.SetWriteDeadline(time.Now().Add(time.Second))
			c.valid()
			go c.handleReceive()
			break
		}

		if waitTime >= WAIT_FRESHOLD_TIME_SEC {
			waitTime = WAIT_FRESHOLD_TIME_SEC
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
	for c.Isvalid() {
		data, err := util.FixedLengthRead(c.conn)
		if err != nil {
			fmt.Println("Wrong to read data.", err)
			continue
		}

		receiver.OnMessage(data)
	}
}
