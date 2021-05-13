package network

import (
	"net"
)

const (
	bitOfSize      = 1
	byteOfSize     = bitOfSize * 1024
	megaByteOfSize = byteOfSize * 1024
	maxPacketSize  = megaByteOfSize * 10
)

type sessionReadFunc func(session *Session, data []byte, size int)
type sessionWriteFunc func(session *Session, bytesTransferred int)
type sessionErrorFunc func(session *Session, err error)
type sessionDisconnected func(session *Session)
type parsePacketHeaderFunc func(conn net.Conn, maxRecvBuffSize int) (int, error)
type buildPacketFunc func(data []byte) []byte

type SessionSettings struct {
	MaxRecvBuffSize int
	MaxSendBuffSize int

	OnRead              sessionReadFunc
	OnWrite             sessionWriteFunc
	OnError             sessionErrorFunc
	OnDisconnected      sessionDisconnected
	OnParsePacketHeader parsePacketHeaderFunc
	OnBuildPacket buildPacketFunc
}

type Session struct {
	socket *Socket
	stop chan struct{}

	maxRecvBuffSize int
	maxSendBuffSize int

	OnRead              sessionReadFunc
	OnWrite             sessionWriteFunc
	OnError             sessionErrorFunc
	OnDisconnected      sessionDisconnected
	OnParsePacketHeader parsePacketHeaderFunc
	OnBuildPacket buildPacketFunc
}

func (session *Session) SetSessionSetting(settings SessionSettings) {
	session.OnRead = settings.OnRead
	session.OnWrite = settings.OnWrite
	session.OnError = settings.OnError
	session.OnDisconnected = settings.OnDisconnected
	session.OnParsePacketHeader = settings.OnParsePacketHeader
	session.OnBuildPacket = settings.OnBuildPacket
	session.maxRecvBuffSize = settings.MaxRecvBuffSize
	session.maxSendBuffSize = settings.MaxSendBuffSize

	if session.OnRead == nil {
		session.OnRead = func(session *Session, data []byte, size int) {
		}
	}

	if session.OnWrite == nil {
		session.OnWrite = func(session *Session, bytesTransferred int) {
		}
	}

	if session.OnError == nil {
		session.OnError = func(session *Session, err error) {
		}
	}

	if session.OnDisconnected == nil {
		session.OnDisconnected = func(session *Session) {
		}
	}

	if session.OnParsePacketHeader == nil {
		session.OnParsePacketHeader = parsePacketHeader
	}

	if session.OnBuildPacket == nil {
		session.OnBuildPacket = buildPacket
	}

	if session.maxRecvBuffSize == 0 {
		session.maxRecvBuffSize = maxPacketSize
	}

	if session.maxSendBuffSize == 0 {
		session.maxSendBuffSize = maxPacketSize
	}
}

func (session *Session) Start() {
	go session.doRecvPacket()
}

func (session *Session) Stop() {
	close(session.stop)
	session.socket.Close()
}

func (session *Session) doRecvPacket() {
	for {
		select {
		case <-session.stop:
			return
		default:
			packetSize, err := session.OnParsePacketHeader(session.socket.conn, session.maxRecvBuffSize)

			if err != nil {
				session.OnError(session, err)
				session.Stop()
				session.OnDisconnected(session)
				return
			}

			packet, err := parsePacketBody(session.socket.conn, packetSize)

			if err != nil {
				session.OnError(session, err)
				session.Stop()
				session.OnDisconnected(session)
				return
			}

			session.OnRead(session, packet, packetSize)
		}
	}
}

func (session *Session) GetMaxRecvBuffSize() int {
	return session.maxRecvBuffSize
}

func (session *Session) GetMaxSendBuffSize() int {
	return session.maxSendBuffSize
}

func (session *Session) SendPacket(data []byte) {
	packet := session.OnBuildPacket(data)
	session.socket.SendPacket(packet)
}

func NewSession(settings SessionSettings, s *Socket) *Session {
	session := &Session{
		socket: s,
	}

	session.SetSessionSetting(settings)

	return session
}
