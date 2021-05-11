package network

import (
	"errors"
	"net"
	"strconv"
)

type acceptorIdGeneratorFunc func() (int64, error)
type acceptorListenFunc func(acceptor *Acceptor)
type acceptorNewSessionFunc func(acceptor *Acceptor, session *Session)
type acceptorErrorFunc func(acceptor *Acceptor, err error)

type AcceptorSettings struct {
	IdGenerator acceptorIdGeneratorFunc

	OnListen     acceptorListenFunc
	OnNewSession acceptorNewSessionFunc
	OnError      acceptorErrorFunc

	SessionSettings SessionSettings
}

type Acceptor struct {
	id int64
	listener net.Listener

	idGenerator acceptorIdGeneratorFunc

	onListen     acceptorListenFunc
	onNewSession acceptorNewSessionFunc
	onError      acceptorErrorFunc

	sessionSettings SessionSettings
}

func (acceptor *Acceptor) SetAcceptorSettings(settings AcceptorSettings) {
	acceptor.idGenerator = settings.IdGenerator
	acceptor.onListen = settings.OnListen
	acceptor.onNewSession = settings.OnNewSession
	acceptor.onError = settings.OnError
	acceptor.sessionSettings = settings.SessionSettings

	if acceptor.idGenerator == nil {
		acceptor.idGenerator = GenerateId()
	}

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
	address := host + ":" + strconv.Itoa(port)
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
		conn, err := acceptor.listener.Accept()

		if err != nil {
			acceptor.onError(acceptor, err)
			acceptor.listener.Close()
			return
		}

		id, err := acceptor.idGenerator()

		if err != nil {
			acceptor.onError(acceptor, errors.New(""))
			acceptor.listener.Close()
			return
		}

		socket := NewSocket(conn)
		session := NewSession(acceptor.sessionSettings, id, socket)
		acceptor.onNewSession(acceptor, session)
		session.Start()
	}
}

func (acceptor *Acceptor) Stop() {
	acceptor.listener.Close()
}

func (acceptor *Acceptor) SetId(id int64) {
	acceptor.id = id
}

func (acceptor *Acceptor) GetId() int64 {
	return acceptor.id
}

func GenerateId() func() (int64, error) {
	var id int64 = 0

	return func() (int64, error) {
		id = id + 1

		return id, nil
	}
}
