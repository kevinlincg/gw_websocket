package gwiface

import (
	"context"
	"net"

	"github.com/gorilla/websocket"
)

/**
連接的定義
client連上時會產生一個新的連接
*/
type Connection interface {
	Start()                                  // 啟動Connection，讓此Connection開始工作
	Stop()                                   // 停止Connection，把此Connection結束
	Context() context.Context                // 回傳ctx，用在自己定義的go routine要把Connection變更狀態
	GetConnection() *websocket.Conn          // 得到內部的socket Conn
	GetConnID() int64                        // 取得ConnectionID
	RemoteAddr() net.Addr                    // 取得Client的IP
	SendMsg(msgID uint32, data []byte) error // 送Message給Client

	SetProperty(key string, value interface{})   //設定一個屬性
	GetProperty(key string) (interface{}, error) //取得一個属性
	RemoveProperty(key string)                   //移除属性

	SetPing()      // 當ping客戶端有回pong時，可以用setping設成有回應
	GetPing() bool // 檢查是否有在時間內取得心跳
	RemovePing()   // 把心跳flag設成false

	IsHeartbeatTimeout() //會一直自己循環的檢查心跳的function
}
