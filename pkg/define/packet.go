package define

import (
	"fmt"
)

type Packet interface {
	GetRawPacket() any
	GetKey() string
	GetDestId() int
	Info() string
	Verbose() string
}

type WrappedPacket struct {
	Packet any
	Key    string
	DestId int
}

func (p WrappedPacket) GetRawPacket() any {
	return p.Packet
}

func (p WrappedPacket) GetKey() string {
	return p.Key
}

func (p WrappedPacket) GetDestId() int {
	return p.DestId
}

func (p WrappedPacket) Info() string {
	return fmt.Sprintf("Key %s DestId %d", p.Key, p.DestId)
}

func (p WrappedPacket) Verbose() string {
	return fmt.Sprintf("Packet %v Key %s DestId %d", p.Packet, p.Key, p.DestId)
}
