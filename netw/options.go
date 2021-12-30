package netw

import "github.com/kevinlincg/gw_websocket/iface"

type Option func(s *Server)

func WithPacket(pack iface.Packet) Option {
	return func(s *Server) {
		s.packet = pack
	}
}
