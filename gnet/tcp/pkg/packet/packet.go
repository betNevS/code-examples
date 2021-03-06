package packet

import (
	"bytes"
	"fmt"
)

const (
	CommandConn = iota + 0x01
	CommandSubmit
)

const (
	CommandConnAck = iota + 0x80
	CommandSubmitAck
)

type Packet interface {
	Decode([]byte) error
	Encode() ([]byte, error)
}

type PktHdr struct {
	CommandID uint8
}

type Submit struct {
	ID      string
	Payload []byte
}

func (s *Submit) Decode(packet []byte) error {
	s.ID = string(packet[:8])
	s.Payload = packet[8:]
	return nil
}

func (s *Submit) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), s.Payload}, nil), nil
}

type SubmitAck struct {
	ID     string
	Result uint8
}

func (s *SubmitAck) Decode(packet []byte) error {
	s.ID = string(packet[:8])
	s.Result = uint8(packet[8])
	return nil
}

func (s *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), []byte{s.Result}}, nil), nil
}

func Decode(packet []byte) (Packet, error) {
	commandID := packet[0]
	switch commandID {
	case CommandSubmit:
		s := Submit{}
		err := s.Decode(packet[1:])
		if err != nil {
			return nil, err
		}
		return &s, nil
	case CommandSubmitAck:
		s := SubmitAck{}
		err := s.Decode(packet[1:])
		if err != nil {
			return nil, err
		}
		return &s, nil
	default:
		return nil, fmt.Errorf("unknown commandID [%d]", commandID)
	}
}

func Encode(p Packet) ([]byte, error) {
	var (
		commandID uint8
		body      []byte
		err       error
	)
	switch t := p.(type) {
	case *Submit:
		commandID = CommandSubmit
		body, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		commandID = CommandSubmitAck
		body, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type [%s]", t)
	}
	return bytes.Join([][]byte{[]byte{commandID}, body}, nil), nil
}
