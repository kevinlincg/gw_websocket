package gwnet

//Message 消息
// 跟前端實際溝通的內容，msgId用來當type，data是任意struct封裝後的資料
type Message struct {
	ID   uint32 `json:"msgId"`
	Data []byte `json:"data"`
}

//NewMsgPackage 建立一個Message物件
func NewMsgPackage(ID uint32, data []byte) *Message {
	return &Message{
		ID:   ID,
		Data: data,
	}
}

//GetMsgID 取得MessageID
func (msg *Message) GetMsgID() uint32 {
	return msg.ID
}

//GetData 取得MessageData
func (msg *Message) GetData() []byte {
	return msg.Data
}

//SetMsgID 變更一個Message的ID
func (msg *Message) SetMsgID(msgID uint32) {
	msg.ID = msgID
}

//SetData 變更一個Message的Data
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
