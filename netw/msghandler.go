package netw

import (
	"strconv"

	"github.com/kevinlincg/gw_websocket/iface"

	"go.uber.org/zap"
)

// MsgHandle -
type MsgHandle struct {
	Apis           map[uint32]iface.Router // 每個MsgID對應的處理方法
	WorkerPoolSize uint32
	TaskQueue      []chan iface.Request // Worker負責取任務的Queue
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]iface.Router),
		WorkerPoolSize: config.WorkerPoolSize,
		// 一個Worker對應一個Queue
		TaskQueue: make([]chan iface.Request, config.WorkerPoolSize),
	}
}

func (mh MsgHandle) DoMsgHandler(request iface.Request) {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error("Call err: ", err)
		}
	}()
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		zap.S().Error("api msgID = ", request.GetMsgID(), " is not FOUND!")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh MsgHandle) AddRouter(msgID uint32, router iface.Router) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api , msgID = " + strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID] = router
}

func (mh MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan iface.Request, 1)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
func (mh MsgHandle) SendMsgToTaskQueue(request iface.Request) {
	workerID := uint32(request.GetConnection().GetConnID()) % mh.WorkerPoolSize
	mh.TaskQueue[workerID] <- request
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan iface.Request) {
	zap.S().Debug("Worker ID = ", workerID, " is started.")
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}
