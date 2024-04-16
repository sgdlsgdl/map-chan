package facade

import (
	"github.com/sgdlsgdl/map-chan/config/concurrency"
	"github.com/sgdlsgdl/map-chan/internal/service/concurrency"
	"github.com/sgdlsgdl/map-chan/pkg/define"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"go.uber.org/zap"
)

type ConcurrencyInterface interface {
	Exec(p define.Packet) error
	ChangeChan(chanSize, chanNum int)
}

func InitConcurrency(options ...concurrencyCfg.Option) ConcurrencyInterface {
	cfg := concurrencyCfg.InitConfig(options...)
	logger.L = logger.WrappedLogger{Logger: cfg.Logger}
	logger.L.Info("InitConcurrency", zap.Any("cfg", cfg))
	return concurrency.NewService(*cfg.ConcurrencyConfig, cfg.ConcurrencyFn)
}
