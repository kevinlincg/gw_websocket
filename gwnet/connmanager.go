package gwnet

import (
	"errors"
	"github.com/kevinlincg/gw_websocket/gwiface"
	"sync"
)

type ConnManager struct {
	connections sync.Map
}

func NewConnManager() *ConnManager {
	return &ConnManager{}
}

func (connMgr *ConnManager) Add(conn gwiface.Connection) {
	connMgr.connections.Store(conn.GetConnID(), conn)
}

func (connMgr *ConnManager) Remove(conn gwiface.Connection) {
	connMgr.connections.Delete(conn.GetConnID())
}

func (connMgr *ConnManager) Get(connID int64) (gwiface.Connection, error) {
	value, ok := connMgr.connections.Load(connID)
	if !ok {
		return nil, errors.New("connection not found")
	}
	return value.(gwiface.Connection), nil
}

func (connMgr *ConnManager) Len() int {
	var length int
	connMgr.connections.Range(func(k, v interface{}) bool {
		length++
		return true
	})
	return length
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connections.Range(func(key, iConn interface{}) bool {
		iConn.(gwiface.Connection).Stop()
		connMgr.connections.Delete(key)
		return true
	})
}

func (connMgr *ConnManager) Search(s gwiface.Search) {
	connMgr.connections.Range(func(_, iConn interface{}) bool {
		if conn := iConn.(gwiface.Connection); conn != nil {
			func() {
				s(conn)
			}()
		}
		return true // Continue to next item
	})
}

func (connMgr *ConnManager) ClearOneConn(connID int64) {
	value, ok := connMgr.connections.Load(connID)
	if ok {
		conn := value.(gwiface.Connection)
		conn.Stop()
		connMgr.connections.Delete(connID)
	}
}
