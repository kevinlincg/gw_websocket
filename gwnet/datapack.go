package gwnet

import (
	"encoding/json"

	"github.com/kevinlincg/gw_websocket/gwiface"
)

type DataPack struct{}

func NewDataPack() gwiface.Packet {
	return &DataPack{}
}

func (dp *DataPack) Pack(msg gwiface.Message) ([]byte, error) {
	msgN := Message{msg.GetMsgID(), msg.GetData()}
	return json.Marshal(&msgN)
}

func (dp *DataPack) Unpack(binaryData []byte) (gwiface.Message, error) {
	var msg Message
	json.Unmarshal(binaryData, &msg)
	return &msg, nil
}
