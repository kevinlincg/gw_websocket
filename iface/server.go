package iface

import (
	"github.com/gin-gonic/gin"
)

/**
定義Server的一個interface
*/
type Server interface {
	Start(c *gin.Context)                  // 啟動Server
	Stop()                                 // 停止Server
	Serve(c *gin.Context)                  // 開始一個服務
	AddRouter(msgID uint32, router Router) // 路由功能：註冊處理Msg的方法

	GetConnMgr() ConnManager // 取得連接管理器，可以從中獲得連結

	SetOnConnStart(func(Connection)) // 設定Server的連接建立時Hook的function
	SetOnConnStop(func(Connection))  // 設定Server的連接中斷時的Hook

	CallOnConnStart(conn Connection) // Call上面SetOnConnStart設定的那個 Hook Function
	CallOnConnStop(conn Connection)  // Call上面SetOnConnStop設定的那個 Hook Function

	Packet() Packet
}
