package tcpnode

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"linkServer/common/util"
	"linkServer/packet"
)

//WaitFreshHoldTimeSec is the longest waiting time when reconnect.
const WaitFreshHoldTimeSec = 600

//NodeClient will build tcp connection to other server.
type NodeClient struct {
	receiver MessageReceiver

	maxMessageLength uint32
	address          string
	conn             net.Conn

	mu        sync.Mutex
	connected bool // protected by mu
	canceled  bool // protected by mu
}

//NewNodeClient return a new NodeClient
func NewNodeClient(address string, receiver MessageReceiver) *NodeClient {
	client := new(NodeClient)

	client.receiver = receiver
	client.address = address
	client.connected = false
	client.canceled = false

	client.connect()

	return client
}

//IsValid judge the tcp connection is ok or not.
func (c *NodeClient) IsValid() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

func (c *NodeClient) valid() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connected = true
}

func (c *NodeClient) unValid() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connected = false
}

//IsCancel judge the NodeClient is cancel or not.
func (c *NodeClient) IsCancel() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.canceled
}

//Cancel the NodeClient.
func (c *NodeClient) Cancel() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.canceled = true
	c.connected = false
	c.conn.Close()
}

//Reset will rebuild the connection.
func (c *NodeClient) Reset() {
	c.unValid()
	if c.conn != nil { // close old connection
		c.conn.Close()
		fmt.Println("We are closing the old connection and begin to reconnect for address : ", c.address)
	}
	c.connect()
}

/**
* This method will block util the connection is built or the client is canceled.
 */
func (c *NodeClient) connect() {

	waitTime := 1
	var err error
	for !c.IsValid() && !c.IsCancel() {
		c.conn, err = net.Dial("tcp", c.address)
		if err == nil {
			c.conn.SetWriteDeadline(time.Now().Add(time.Second))
			c.valid()
			go c.handleReceive()
			break
		}

		if waitTime >= WaitFreshHoldTimeSec {
			waitTime = WaitFreshHoldTimeSec
		} else {
			waitTime = 2 * waitTime
		}

		fmt.Println("Connect failed. We are now going to sleep and then retry. waitTime is ", waitTime)
		time.Sleep(waitTime * time.Second)
	}

	return
}

//Send is to send data to connection.
func (c *NodeClient) Send(data []byte) error {
	datalen := len(data)
	if datalen > int(c.maxMessageLength) {
		return errors.New("data too big")
	}

	lenBuf := packet.Uint32ToBytes(uint32(datalen))
	data = append(lenBuf, data...)
	_, err := c.conn.Write(data)
	if err != nil {
		go c.Reset()
	}

	return err
}

func (c *NodeClient) handleReceive() {
	for c.IsValid() {
		data, err := util.FixedLengthRead(c.conn)
		if err != nil {
			fmt.Println("Wrong to read data.", err)
			continue
		}

		c.receiver.OnMessage(data)
	}
}
