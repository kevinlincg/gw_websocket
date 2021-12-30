package netw

import "github.com/kevinlincg/gw_websocket/iface"

var (
	config *iface.Config
)

func SetConfig(c *iface.Config) {
	config = c
}
