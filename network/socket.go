package network

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
)

type Socket struct {
	conn net.Conn

	LocalHost    string
	LocalPort    int
	LocalAddress string
	RemoteHost   string
	RemotePort   int
}

func (s *Socket) Connect(host string, port int) bool {
	address := ComposeAddressByHostAndPort(host, port)
	conn, err := net.Dial("tcp", address)

	if err != nil {
		return false
	}

	localHost, localPort, err := SplitHostAndPort(conn.LocalAddr().String())

	if err != nil {
		return false
	}

	s.conn = conn
	s.LocalHost = localHost
	s.LocalPort = localPort
	s.RemoteHost = host
	s.RemotePort = port

	return true
}

func (s *Socket) GetLocalHost() string {
	return s.LocalHost
}

func (s *Socket) GetLocalPort() int {
	return s.LocalPort
}

func (s *Socket) GetLocalAddress() string {
	return s.conn.LocalAddr().String()
}

func (s *Socket) GetRemoteHost() string {
	return s.RemoteHost
}

func (s *Socket) GetRemotePort() int {
	return s.RemotePort
}

func (s *Socket) GetRemoteAddress() string {
	return s.conn.RemoteAddr().String()
}

func (s *Socket) Close() error {
	return s.conn.Close()
}

func (s *Socket) ReadSome(size int) ([]byte, error) {
	packet := make([]byte, size)
	n, err := io.ReadFull(s.conn, packet)

	if err != nil {
		return nil, err
	}

	if n != len(packet) {
		return nil, errors.New("")
	}

	return packet, nil
}

func (s *Socket) SendPacket(data []byte) (int, error) {
	size, err := s.conn.Write(data)

	if err != nil {
		return size, err
	}

	if size != len(data) {
		return size, errors.New("")
	}

	return size, nil
}

func parsePacketHeader(conn net.Conn, maxRecvBuffSize int) (int, error) {
	packetHeader := make([]byte, 4)
	size, err := io.ReadFull(conn, packetHeader)

	if err != nil {
		return 0, err
	}

	if size != 4 {
		return size, errors.New("")
	}

	return int(binary.BigEndian.Uint32(packetHeader)), nil
}

func parsePacketBody(conn net.Conn, packetSize int) ([]byte, error) {
	packet := make([]byte, packetSize)
	size, err := io.ReadFull(conn, packet)

	if err != nil {
		return nil, err
	}

	if size != packetSize {
		return nil, errors.New("")
	}

	return packet, nil
}

func buildPacket(data []byte) []byte {
	packetSize := len(data)
	totalPacketSize := 4 + packetSize
	packet := make([]byte, totalPacketSize)
	binary.BigEndian.PutUint32(packet, uint32(packetSize))
	copy(packet[4:], data)

	return packet
}

func NewSocket(conn net.Conn) *Socket {
	s := &Socket{}

	localHost, localPort, err := SplitHostAndPort(conn.LocalAddr().String())

	if err != nil {
		return nil
	}

	remoteHost, remotePort, err := SplitHostAndPort(conn.RemoteAddr().String())

	if err != nil {
		return nil
	}

	s.conn = conn
	s.LocalHost = localHost
	s.LocalPort = localPort
	s.RemoteHost = remoteHost
	s.RemotePort = remotePort

	return s
}

func ComposeAddressByHostAndPort(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}

func SplitHostAndPort(address string) (string, int, error) {
	addressParts := strings.Split(address, ":")
	count := len(addressParts)

	if count > 2 {
		return "", 0, errors.New("too many parameters in address")
	} else if count < 2 {
		return "", 0, errors.New("missing port in address")
	}

	port, _:= strconv.Atoi(addressParts[1])

	return addressParts[0], port, nil
}