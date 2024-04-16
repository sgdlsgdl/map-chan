package poolCfg

import (
	"github.com/sgdlsgdl/map-chan/pkg/connector"
	"github.com/sgdlsgdl/map-chan/pkg/discovery"
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

func WithPoolSendChanTimeout(sendChanTimeout time.Duration) Option {
	if sendChanTimeout <= 0 {
		sendChanTimeout = map_chan.DefaultSendChanTimeout
	}
	return func(c *Config) {
		c.PoolConfig.SendChanTimeout = sendChanTimeout
	}
}

func WithPoolChanSize(chanSize int) Option {
	if chanSize <= 0 {
		chanSize = map_chan.DefaultChanSize
	}
	return func(c *Config) {
		c.PoolConfig.ChanSize = chanSize
	}
}

func WithDiscovery(discovery discovery.Discovery) Option {
	return func(c *Config) {
		c.PoolConfig.Discovery = discovery
	}
}

func WithConnectorHandler(connectorHandler connector.Handler) Option {
	return func(c *Config) {
		c.PoolConfig.ConnectorHandler = connectorHandler
	}
}
