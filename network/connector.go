package network

import (
	"net"
)

type connectorConnectedFunc func(connector *Connector)
type connectorDisConnected func(connector *Connector, session *Session)
type connectorErrorFunc func(connector *Connector, err error)

type ConnectorSettings struct {
	OnConnected connectorConnectedFunc
	OnDisconnected connectorDisConnected
	OnError connectorErrorFunc

	SessionSettings SessionSettings
}

type Connector struct {
	session *Session

	onConnected connectorConnectedFunc
	onDisconnected connectorDisConnected
	onError connectorErrorFunc

	sessionSettings SessionSettings
}

func (connector *Connector) SetConnectorSettings(settings ConnectorSettings) {
	connector.onConnected = settings.OnConnected
	connector.onDisconnected = settings.OnDisconnected
	connector.onError = settings.OnError
	connector.sessionSettings = settings.SessionSettings

	if connector.onConnected == nil {
		connector.onConnected = func(connector *Connector) {
		}
	}

	if connector.onDisconnected == nil {
		connector.onDisconnected = func(connector *Connector, session *Session) {
		}
	}

	if connector.onError == nil {
		connector.onError = func(connector *Connector, err error) {
		}
	}
}

func (connector *Connector) Connect(host string, port int) bool {
	address := ComposeAddressByHostAndPort(host, port)
	conn, err := net.Dial("tcp", address)

	if err != nil {
		connector.onError(connector, err)
		return false
	}

	s := NewSocket(conn)
	session := NewSession(connector.sessionSettings, s)
	connector.session = session

	return true
}

func (connector *Connector) Start() {
	connector.session.doRecvPacket()
}

func (connector *Connector) Stop() {
	connector.session.Stop()
}

func NewConnector(settings ConnectorSettings) *Connector {
	connector := &Connector{}
	connector.SetConnectorSettings(settings)

	return connector
}
