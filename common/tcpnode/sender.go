package tcpnode



//Send message to another node
type MessageSender interface {

	/*
	* this method will block util the data sent
	* or face an error.
	* this method will send more first 4 bytes which
	* indicate the length of data.
	*/
	func Send(data []byte) error

}
