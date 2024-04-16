package discovery

import "sync"

type ChangeEvent string

const (
	ChangeEventDel  ChangeEvent = "Del"
	ChangeEventAddr ChangeEvent = "Addr"
	ChangeEventNum  ChangeEvent = "Num"
)

type Discovery interface {
	GetDest(destId int) (Dest, bool)
	CloneMap() map[int]Dest
	UpdateAndDiff(Discovery) map[ChangeEvent][]Dest
	sync.Locker
}

type Dest interface {
	GetDestId() int
	GetName() string
	GetAddr() string
	GetConnNum() int
	SetConnNum(int) Dest
}
