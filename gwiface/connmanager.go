package gwiface

/*
	Connection的管理器 interface
*/
type Search func(Connection)
type ConnManager interface {
	Add(conn Connection)                  // 增加一個conn
	Remove(conn Connection)               // 移除conn
	Get(connID int64) (Connection, error) // 取得Conn使用ConnID
	Len() int                             // 取得總Conn數
	Search(Search)                        // 尋找Conn
	ClearConn()                           // 刪除並且Stop所有Conn
}
