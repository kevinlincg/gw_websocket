package netw

import "github.com/kevinlincg/gw_websocket/iface"

//Request 請求
type Request struct {
	conn iface.Connection //已經跟client建立好的連線
	msg  iface.Message    //client端請求的資料
}

//GetConnection 取得connection
func (r *Request) GetConnection() iface.Connection {
	return r.conn
}

//GetData 取得資料
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//GetMsgID 取得msgID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
