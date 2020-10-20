package mcping

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
)

// Pong is the server's response to a Ping request
type Pong struct {
	Version struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int64 `json:"max"`
		Online int64 `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	}
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	Favicon string `json:"favicon"`
}

// we dont care about errors here
func do(_ error)         {}
func don(_ int, _ error) {}

func sendHandshake(ctx context.Context, wire *bufio.ReadWriter, hostport Server) error {
	// packet id
	do(wire.WriteByte(0x00))

	// protocol version
	do(wire.WriteByte(math.MaxUint8 - 1))

	// server host
	varint := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(varint, uint64(len(hostport.Host)))
	don(wire.Write(varint[:n]))
	don(wire.WriteString(hostport.Host))

	// server port
	do(binary.Write(wire, binary.BigEndian, hostport.Port))

	// status
	do(wire.WriteByte(0x01))
	return wire.Flush()
}

func sendPing(ctx context.Context, wire *bufio.ReadWriter) error {
	do(wire.WriteByte(0x00))
	return wire.Flush()
}

func parsePong(ctx context.Context, wire *bufio.ReadWriter) (*Pong, error) {
	l, err := binary.ReadUvarint(wire)
	if err != nil {
		return nil, fmt.Errorf("could not read packet length: %s", err.Error())
	}

	p := make([]byte, l)
	_, err = io.ReadFull(wire, p)
	if err != nil {
		return nil, fmt.Errorf("could not read packet: %s", err.Error())
	}

	// packet id
	_, n1 := binary.Uvarint(p)
	if n1 <= 0 {
		return nil, errors.New("could not read packet id")
	}

	// string varint
	_, n2 := binary.Uvarint(p)
	if n2 <= 0 {
		return nil, errors.New("could not read string varint")
	}

	var pong Pong
	err = json.Unmarshal(p[n1+n2:], &pong)
	if err != nil {
		return nil, err
	}
	return &pong, nil
}

// Ping pings a minecraft server
func Ping(ctx context.Context, conn net.Conn, hostport Server) (*Pong, error) {
	wire := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	err := sendHandshake(ctx, wire, hostport)
	if err != nil {
		return nil, fmt.Errorf("error sending handshake: %s", err.Error())
	}

	err = sendPing(ctx, wire)
	if err != nil {
		return nil, fmt.Errorf("error sending ping: %s", err.Error())
	}

	pong, err := parsePong(ctx, wire)
	if err != nil {
		return nil, fmt.Errorf("error parsing pong: %s", err.Error())
	}
	return pong, nil
}
