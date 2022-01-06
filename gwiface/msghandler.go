package gwiface

/*
	Msg管理器的封裝
*/
type MsgHandle interface {
	DoMsgHandler(request Request)          // 馬上用non-blocking的方式處理msg
	AddRouter(msgID uint32, router Router) // 增加一個msg的處理邏輯
	StartWorkerPool()                      // 啟動worker工作池
	SendMsgToTaskQueue(request Request)    // msg給TaskQueue,由worker進行處理
}
