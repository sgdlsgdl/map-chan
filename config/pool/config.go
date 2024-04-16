package poolCfg

import (
	"github.com/sgdlsgdl/map-chan/pkg/connector"
	"github.com/sgdlsgdl/map-chan/pkg/discovery"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"github.com/sgdlsgdl/map-chan/pkg/map_chan"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Config struct {
	Logger logger.LoggerInterface

	PoolConfig *PoolConfig `json:"pool_config"`
}

type PoolConfig struct {
	SendChanTimeout  time.Duration       `json:"send_chan_timeout"`
	ChanSize         int                 `json:"chan_size"`
	Discovery        discovery.Discovery `json:"-"`
	ConnectorHandler connector.Handler   `json:"-"`
}

func NewDefaultConfig() *Config {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel)
	lo := zap.New(core)
	return &Config{
		Logger: lo,
		PoolConfig: &PoolConfig{
			SendChanTimeout: map_chan.DefaultSendChanTimeout,
			ChanSize:        map_chan.DefaultChanSize,
		},
	}
}

func InitConfig(options ...Option) *Config {
	cfg := NewDefaultConfig()
	for _, option := range options {
		option(cfg)
	}
	return cfg
}
