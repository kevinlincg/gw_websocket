package gwnet

import (
	"errors"
	"github.com/kevinlincg/gw_websocket/gwiface"
	"sync"
)

type ConnManager struct {
	connections map[int64]gwiface.Connection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[int64]gwiface.Connection),
	}
}

func (connMgr *ConnManager) Add(conn gwiface.Connection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	connMgr.connections[conn.GetConnID()] = conn
}

func (connMgr *ConnManager) Remove(conn gwiface.Connection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	delete(connMgr.connections, conn.GetConnID())
}

func (connMgr *ConnManager) Get(connID int64) (gwiface.Connection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

func (connMgr *ConnManager) Len() int {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	length := len(connMgr.connections)
	return length
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}
}

func (connMgr *ConnManager) Search(s gwiface.Search) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	for _, conn := range connMgr.connections {
		s(conn)
	}
}

func (connMgr *ConnManager) ClearOneConn(connID int64) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	connections := connMgr.connections
	if conn, ok := connections[connID]; ok {
		conn.Stop()
		delete(connections, connID)
		return
	}
	return
}
