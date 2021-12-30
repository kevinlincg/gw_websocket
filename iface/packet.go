package iface

/*
	封包跟拆包封包的資料用的
	可以實作這個來增加自己的加密
*/
type Packet interface {
	Pack(msg Message) ([]byte, error)
	Unpack([]byte) (Message, error)
}
