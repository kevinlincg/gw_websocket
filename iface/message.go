package iface

/*
	把Request的Message封裝到這個Message裡面用的
*/
type Message interface {
	GetMsgID() uint32
	GetData() []byte

	SetMsgID(uint32)
	SetData([]byte)
}
