package discovery

import (
	"sync"
)

var (
	defaultChanNum = 16
)

type DestMap struct {
	DestIdToDest map[int]Dest
	mutex        sync.RWMutex
}

func (d *DestMap) Lock() {
	d.mutex.Lock()
}

func (d *DestMap) Unlock() {
	d.mutex.Unlock()
}

func (d *DestMap) GetDest(destId int) (Dest, bool) {
	d.mutex.RLock()
	dest, exist := d.DestIdToDest[destId]
	d.mutex.RUnlock()
	if exist && dest.GetConnNum() <= 0 {
		dest = dest.SetConnNum(defaultChanNum)
	}
	return dest, exist
}

func (d *DestMap) CloneMap() map[int]Dest {
	m := make(map[int]Dest)
	d.mutex.RLock()
	for destId, dest := range d.DestIdToDest {
		m[destId] = dest
	}
	d.mutex.RUnlock()
	return m
}

func (d *DestMap) UpdateAndDiff(newD Discovery) map[ChangeEvent][]Dest {
	var changeDelList, changeAddrList, changeNumList []Dest
	newM := newD.CloneMap()
	d.Lock()
	for destId, dest := range d.DestIdToDest {
		newDest, ok := newM[destId]
		if !ok {
			changeDelList = append(changeDelList, dest)
			continue
		}
		if dest.GetAddr() != newDest.GetAddr() {
			changeAddrList = append(changeAddrList, newDest)
		}
		if dest.GetConnNum() != newDest.GetConnNum() {
			changeNumList = append(changeNumList, newDest)
		}
	}
	d.DestIdToDest = newM
	d.Unlock()

	res := make(map[ChangeEvent][]Dest)
	res[ChangeEventDel] = changeDelList
	res[ChangeEventAddr] = changeAddrList
	res[ChangeEventNum] = changeNumList
	return res
}

type DestServer struct {
	DestId  int    `json:"dest_id"`
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	ConnNum int    `json:"conn_num"`
}

func (d DestServer) GetDestId() int {
	return d.DestId
}

func (d DestServer) GetName() string {
	return d.Name
}

func (d DestServer) GetAddr() string {
	return d.Addr
}

func (d DestServer) GetConnNum() int {
	return d.ConnNum
}

func (d DestServer) SetConnNum(i int) Dest {
	d.ConnNum = i
	return d
}
