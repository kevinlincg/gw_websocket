package gwnet

import (
	"github.com/kevinlincg/gw_websocket/gwiface"
	"net/http"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	GlobalServer gwiface.Server
)

// Server interface的實現，定義一個Server的類型
type Server struct {
	sesIDGen int64 // 紀錄最新的ConnID，ConnID用流水號產生

	msgHandler gwiface.MsgHandle

	ConnMgr gwiface.ConnManager

	OnConnStart func(conn gwiface.Connection)
	OnConnStop  func(conn gwiface.Connection)

	packet gwiface.Packet
}

// NewServer 建立一個Server
func NewServer(opt ...Option) gwiface.Server {
	s := &Server{
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		packet:     NewDataPack(),
	}
	for _, option := range opt {
		option(s)
	}
	s.msgHandler.StartWorkerPool()
	return s
}

// ============== 實現 iface.Server 的function ========

func (s *Server) Start(c *gin.Context) {
	// 等待客户端建立连接请求
	var (
		err      error
		wsSocket *websocket.Conn
	)
	if wsSocket, err = Upgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
		return
	}
	if s.ConnMgr.Len() >= config.MaxConn {
		_ = wsSocket.Close()
		return
	}
	dealConn := NewConnection(s, wsSocket, atomic.AddInt64(&s.sesIDGen, 1), s.msgHandler)
	dealConn.Start()
}

func (s *Server) Stop() {
	zap.S().Info("[STOP] server...")
	// 將其他需要Clear的東西，都要一起清理掉
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve(c *gin.Context) {
	s.Start(c)
	select {}
}

func (s *Server) AddRouter(msgID uint32, router gwiface.Router) {
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnMgr() gwiface.ConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(gwiface.Connection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(gwiface.Connection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn gwiface.Connection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn gwiface.Connection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}

func (s *Server) Packet() gwiface.Packet {
	return s.packet
}
