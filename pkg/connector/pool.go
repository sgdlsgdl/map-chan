package connector

import "github.com/sgdlsgdl/map-chan/pkg/define"

type Pool interface {
	Get() (Conn, error)
	Close()
	Len() int
}

type Conn interface {
	Write(define.Packet) error
	IsClosed() bool
	Close() error
	Destroy()
	Info() string
}
