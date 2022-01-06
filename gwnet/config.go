package gwnet

import "github.com/kevinlincg/gw_websocket/gwiface"

var (
	config *gwiface.Config
)

func SetConfig(c *gwiface.Config) {
	config = c
}
