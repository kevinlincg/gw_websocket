package iface

/*
	Request 接口：
	把client端的Conn跟Msg封裝在一起成為一個Request

	主邏輯處理事處理Request，ws底層是讀Msg出來
*/
type Request interface {
	GetConnection() Connection
	GetData() []byte
	GetMsgID() uint32
}
