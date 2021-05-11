package network

import (
	"net"
	"strconv"
)

type connectorConnectedFunc func(connector *Connector)
type connectorErrorFunc func(connector *Connector, err error)

type ConnectorSettings struct {
	OnConnected connectorConnectedFunc
	OnError     connectorErrorFunc

	SessionSettings SessionSettings
}

type Connector struct {
	id int64
	session *Session

	onConnected connectorConnectedFunc
	onError     connectorErrorFunc

	sessionSettings SessionSettings
}

func (connector *Connector) SetConnectorSettings(settings ConnectorSettings) {
	connector.onConnected = settings.OnConnected
	connector.onError = settings.OnError
	connector.sessionSettings = settings.SessionSettings

	if connector.onConnected == nil {
		connector.onConnected = func(connector *Connector) {
		}
	}

	if connector.onError == nil {
		connector.onError = func(connector *Connector, err error) {
		}
	}
}

func (connector *Connector) Connect(host string, port int) bool {
	address := host + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", address)

	if err != nil {
		connector.onError(connector, err)
		return false
	}

	s := NewSocket(conn)
	session := NewSession(connector.sessionSettings, 0, s)
	connector.session = session

	return true
}

func (connector *Connector) Start() {
	connector.doRecvPacket()
}

func (connector *Connector) doRecvPacket() {
	connector.session.Start()
}

func (connector *Connector) SetId(id int64) {
	connector.id = id
}

func (connector *Connector) GetId() int64 {
	return connector.id
}

func NewConnector(settings ConnectorSettings) *Connector {
	connector := &Connector{}
	connector.SetConnectorSettings(settings)

	return connector
}
