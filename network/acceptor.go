package network

import (
	"net"
)

type acceptorListenFunc func(acceptor *Acceptor)
type acceptorNewSessionFunc func(acceptor *Acceptor, session *Session)
type acceptorErrorFunc func(acceptor *Acceptor, err error)

type AcceptorSettings struct {
	OnListen     acceptorListenFunc
	OnNewSession acceptorNewSessionFunc
	OnError      acceptorErrorFunc

	SessionSettings SessionSettings
}

type Acceptor struct {
	stop chan struct{}
	listener net.Listener

	onListen     acceptorListenFunc
	onNewSession acceptorNewSessionFunc
	onError      acceptorErrorFunc

	sessionSettings SessionSettings
}

func (acceptor *Acceptor) SetAcceptorSettings(settings AcceptorSettings) {
	acceptor.onListen = settings.OnListen
	acceptor.onNewSession = settings.OnNewSession
	acceptor.onError = settings.OnError
	acceptor.sessionSettings = settings.SessionSettings

	if acceptor.onListen == nil {
		acceptor.onListen = func(acceptor *Acceptor) {
		}
	}

	if acceptor.onNewSession == nil {
		acceptor.onNewSession = func(acceptor *Acceptor, session *Session) {
		}
	}

	if acceptor.onError == nil {
		acceptor.onError = func(acceptor *Acceptor, err error) {
		}
	}
}

func (acceptor *Acceptor) Start(host string, port int) bool {
	address := ComposeAddressByHostAndPort(host, port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		acceptor.onError(acceptor, err)
		return false
	}

	acceptor.listener = listener
	go acceptor.doAccept()

	return true
}

func (acceptor *Acceptor) doAccept() {
	for {
		select {
		case <-acceptor.stop:
			return
		default:
			conn, err := acceptor.listener.Accept()

			if err != nil {
				acceptor.onError(acceptor, err)
				acceptor.listener.Close()
				return
			}

			socket := NewSocket(conn)
			session := NewSession(acceptor.sessionSettings, socket)
			acceptor.onNewSession(acceptor, session)
			session.Start()
		}
	}
}

func (acceptor *Acceptor) Stop() {
	close(acceptor.stop)
	acceptor.listener.Close()
}
