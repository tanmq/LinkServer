package tcpnode



type MessageReceiver interface {

	/**
	* when an packet data received, it would call this method.
	* data here is fixed length which means the first four bytes indicate
	* how much the packet is and data exclude that four bytes.
	*/
	OnMessage(data []byte)
}
