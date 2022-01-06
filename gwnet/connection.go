package gwnet

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/kevinlincg/gw_websocket/gwiface"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

// Connection
type Connection struct {
	Server gwiface.Server //這個連線屬於哪一個Server

	Conn *websocket.Conn //本連線的websocket底層結構

	ConnID int64 //ConnectionID可以當成SessionID，應該要唯一的

	MsgHandler gwiface.MsgHandle //管理Msg的處理

	Heartbeat             bool //用來給Timer檢查是否有心跳
	OfflineHeartBeatCount int  //已經累積幾次沒有心跳，超過n次要把他斷線

	ctx    context.Context //用來取消這個Connection用的
	cancel context.CancelFunc

	msgChan      chan []byte //SendMsg給Client時 用這個channel當作中介
	sync.RWMutex             //讓Connection本身也是一個RWMutex

	property     map[string]interface{} //可以用來給Connection放一些資訊
	propertyLock sync.Mutex

	isClosed bool //連線已經關閉

	writeWait time.Duration
}

// NewConnection 建立一個新的Connection，有新用戶連上時都會新建
func NewConnection(s gwiface.Server, conn *websocket.Conn, connID int64, msgHandler gwiface.MsgHandle) *Connection {
	// 初始化Conn属性
	c := &Connection{
		Server:     s,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		Heartbeat:  false,
		msgChan:    make(chan []byte, 1),
		property:   nil,
		writeWait:  time.Duration(config.WriteDeadlineDelay) * time.Second,
	}
	// 把新的Connection加入ConnMgr
	c.Server.GetConnMgr().Add(c)
	c.IsHeartbeatTimeout()

	return c
}

// StartWriter 寫Message出去用的goroutine，主要是送msg給Client
func (c *Connection) StartWriter() {
	zap.S().Debug("StartWriter [Writer Goroutine is running], Connection: ", c.ConnID)
	defer zap.S().Debug(c.RemoteAddr().String(), "[conn Writer exit!]")
	for {
		select {
		case data := <-c.msgChan:
			// data coming
			c.Conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err := c.Conn.WriteMessage(config.MessageType, data); err != nil {
				zap.S().Error("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// StartReader 用來讀Message，取得Message後送給Worker處理
func (c *Connection) StartReader() {
	zap.S().Debug("StartReader [Reader Goroutine is running] , Connection: ", c.ConnID)
	defer zap.S().Debug(c.RemoteAddr().String(), "[conn Reader exit!]")

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 讀取Client的Msg
			t, msgData, err := c.Conn.ReadMessage()
			if err != nil {
				goto Wrr
			}
			if t != config.MessageType {
				c.Stop()
				continue
			}
			// 拆包，得到msgID 和 data 放在msg中
			msg, err := c.Server.Packet().Unpack(msgData)
			if err != nil {
				zap.S().Error("unpack error ", err)
				goto Wrr
			}
			// 把msg轉成Request物件，然後傳給MsgHandler處理
			req := Request{
				conn: c,
				msg:  msg,
			}
			if config.WorkerPoolSize > 0 {
				c.MsgHandler.SendMsgToTaskQueue(&req)
			} else {
				go c.MsgHandler.DoMsgHandler(&req)
			}
		}
	}
Wrr:
	c.Stop()
}

//Start 啟動Connection，讓此Connection開始工作
func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go c.StartWriter()
	c.Server.CallOnConnStart(c)

	// StartReader 一個for loop持續的ReadMessage，所以會卡在這邊直到c.ctx 或c.cancel關閉
	c.StartReader()

}

//Stop 停止Connection，把此Connection結束
func (c *Connection) Stop() {
	c.Lock()
	defer c.Unlock()

	c.Server.CallOnConnStop(c)

	if c.isClosed == true {
		return
	}

	zap.S().Debug("Conn Stop()...ConnID = ", c.ConnID)
	// 關閉Writer
	c.cancel()
	// 關閉socket
	_ = c.Conn.Close()
	// 關閉Reader
	close(c.msgChan)
	// 設置成已經關閉
	c.isClosed = true

	c.Server.GetConnMgr().Remove(c)
}

//Context 回傳ctx，用在自己定義的go routine要把Connection變更狀態
func (c *Connection) Context() context.Context {
	return c.ctx
}

//GetConnection 得到內部的socket Conn
func (c *Connection) GetConnection() *websocket.Conn {
	return c.Conn
}

//GetConnID 取得ConnectionID
func (c *Connection) GetConnID() int64 {
	return c.ConnID
}

//RemoteAddr 取得Client的IP
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//SendMsg 送Message給Client
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	c.RUnlock()

	dp := c.Server.Packet()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		zap.S().Error("pack error msg = ", msgID, data)
		return errors.New("pack error msg ")
	}

	c.msgChan <- msg
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

//SetPing 設定有HeartBeat
func (c *Connection) SetPing() {
	c.Lock()
	c.Heartbeat = true
	c.Unlock()
}

//GetPing 檢查是否有HeartBeat
func (c *Connection) GetPing() bool {
	return c.Heartbeat
}

//RemovePing HeartBeat
func (c *Connection) RemovePing() {
	c.Lock()
	c.Heartbeat = false
	c.Unlock()
}

//IsHeartbeatTimeout 自己會陷入無限循環，直到GetPing = false 會中斷連線
func (c *Connection) IsHeartbeatTimeout() {
	pingTime := time.Second * time.Duration(config.PingTime+1)
	delayFunc := func() {
		if !c.GetPing() {
			c.Stop()
		} else {
			c.RemovePing()
			c.IsHeartbeatTimeout()
		}
	}
	time.AfterFunc(pingTime, delayFunc)
	return
}
