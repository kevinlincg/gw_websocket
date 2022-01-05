package netw

import (
	"bytes"
	"encoding/binary"

	"github.com/kevinlincg/gw_websocket/iface"
)

type DataPack struct{}

func NewDataPack() iface.Packet {
	return &DataPack{}
}

func (dp *DataPack) Pack(msg iface.Message) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (iface.Message, error) {
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.ID); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.Data); err != nil {
		return nil, err
	}
	return msg, nil
}
