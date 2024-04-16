package pool

import "strconv"

const (
	Success        = SendPacketCode(0)
	GetDestErr     = SendPacketCode(-10001)
	GetPoolErr     = SendPacketCode(-10002)
	GetConnErr     = SendPacketCode(-10003)
	WritePacketErr = SendPacketCode(-10004)
)

type SendPacketCode int

func (c SendPacketCode) String() string {
	return strconv.Itoa(int(c))
}
