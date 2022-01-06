package gwnet

import "github.com/kevinlincg/gw_websocket/gwiface"

type Option func(s *Server)

func WithPacket(pack gwiface.Packet) Option {
	return func(s *Server) {
		s.packet = pack
	}
}
