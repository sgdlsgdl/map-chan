package concurrencyCfg

import (
	"context"
	"github.com/sgdlsgdl/map-chan/pkg/define"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"github.com/sgdlsgdl/map-chan/pkg/map_chan"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Config struct {
	Logger logger.LoggerInterface

	ConcurrencyFn     func(ctx context.Context, p define.Packet) `json:"-"`
	ConcurrencyConfig *ConcurrencyConfig                         `json:"concurrency_config"`
}

type ConcurrencyConfig struct {
	SendChanTimeout time.Duration `json:"send_chan_timeout"`
	ChanSize        int           `json:"chan_size"`
	ChanNum         int           `json:"chan_num"`
}

func NewDefaultConfig() *Config {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel)
	lo := zap.New(core)
	return &Config{
		Logger: lo,
		ConcurrencyConfig: &ConcurrencyConfig{
			SendChanTimeout: map_chan.DefaultSendChanTimeout,
			ChanSize:        map_chan.DefaultChanSize,
			ChanNum:         map_chan.DefaultChanNum,
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
