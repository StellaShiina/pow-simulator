package p2p

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

func WriteMsg(conn net.Conn, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))
	if _, err := conn.Write(header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("write body: %w", err)
	}
	return nil
}

func ReadMsg(conn net.Conn) (*Message, error) {
	// Read header
	header := make([]byte, 4)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	length := binary.BigEndian.Uint32(header)
	if length == 0 {
		return nil, fmt.Errorf("zero-length message")
	}

	// Read body
	body := make([]byte, length)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return &msg, nil
}

func SendMessage(conn net.Conn, msgType string, data any) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return WriteMsg(conn, &Message{Type: msgType, Data: raw})
}

func DialWithRetry(addr string, retry int) (net.Conn, error) {
	var conn net.Conn
	var err error
	for range retry {
		conn, err = net.Dial("tcp", addr)
		if err == nil {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("dial %s after %d retries: %w", addr, retry, err)
}
