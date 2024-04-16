package concurrencyCfg

import (
	"context"
	"github.com/sgdlsgdl/map-chan/pkg/define"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"github.com/sgdlsgdl/map-chan/pkg/map_chan"
	"time"
)

type Option func(c *Config)

func WithLogger(logger logger.LoggerInterface) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

func WithConcurrencyFn(fn func(ctx context.Context, p define.Packet)) Option {
	return func(c *Config) {
		c.ConcurrencyFn = fn
	}
}

func WithConcurrencySendChanTimeout(sendChanTimeout time.Duration) Option {
	if sendChanTimeout <= 0 {
		sendChanTimeout = map_chan.DefaultSendChanTimeout
	}
	return func(c *Config) {
		c.ConcurrencyConfig.SendChanTimeout = sendChanTimeout
	}
}

func WithConcurrencyChanSize(chanSize int) Option {
	if chanSize <= 0 {
		chanSize = map_chan.DefaultChanSize
	}
	return func(c *Config) {
		c.ConcurrencyConfig.ChanSize = chanSize
	}
}

func WithConcurrencyChanNum(chanNum int) Option {
	if chanNum <= 0 {
		chanNum = map_chan.DefaultChanNum
	}
	return func(c *Config) {
		c.ConcurrencyConfig.ChanNum = chanNum
	}
}
