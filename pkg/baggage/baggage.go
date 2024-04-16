package baggage

import (
	"context"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
)

type baggageKey struct{}

type Baggage struct {
	val interface{}
}

func NewBaggage() context.Context {
	return context.WithValue(context.Background(), baggageKey{}, &Baggage{})
}

func SetVal(ctx context.Context, val interface{}) {
	bg, ok := ctx.Value(baggageKey{}).(*Baggage)
	if !ok {
		logger.L.Error("SetVal Baggage not found")
		return
	}
	bg.val = val
}

func GetVal(ctx context.Context) interface{} {
	bg, ok := ctx.Value(baggageKey{}).(*Baggage)
	if !ok {
		logger.L.Error("GetVal Baggage not found")
		return nil
	}
	return bg.val
}
