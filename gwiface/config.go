package gwiface

type Config struct {
	PingTime       int    // HeartBeat檢查的時間
	MaxConn        int    // 允許的最大連線數
	WorkerPoolSize uint32 // Working Pool的數量

	// 定義在gorilla/websocket/conn.go內
	// The message types are defined in RFC 6455, section 11.8.
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	//TextMessage = 1
	//
	//// BinaryMessage denotes a binary data message.
	//BinaryMessage = 2
	//
	//// CloseMessage denotes a close control message. The optional message
	//// payload contains a numeric code and text. Use the FormatCloseMessage
	//// function to format a close message payload.
	//CloseMessage = 8
	//
	//// PingMessage denotes a ping control message. The optional message payload
	//// is UTF-8 encoded text.
	//PingMessage = 9
	//
	//// PongMessage denotes a pong control message. The optional message payload
	//// is UTF-8 encoded text.
	//PongMessage = 10
	//一般使用應該用1或2就好
	MessageType int
}
