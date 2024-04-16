package connector

import (
	"github.com/sgdlsgdl/map-chan/pkg/discovery"
)

type Handler interface {
	GetOrInitPool(dest discovery.Dest) (Pool, error)
	BatchClose([]discovery.Dest)
}
