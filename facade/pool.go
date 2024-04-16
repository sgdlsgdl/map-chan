package facade

import (
	"github.com/sgdlsgdl/map-chan/config/pool"
	"github.com/sgdlsgdl/map-chan/internal/service/pool"
	"github.com/sgdlsgdl/map-chan/pkg/define"
	"github.com/sgdlsgdl/map-chan/pkg/discovery"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"go.uber.org/zap"
)

type PoolInterface interface {
	Push(p define.Packet) error
	Reload(discovery discovery.Discovery)
}

func InitPool(options ...poolCfg.Option) PoolInterface {
	cfg := poolCfg.InitConfig(options...)
	logger.L = logger.WrappedLogger{Logger: cfg.Logger}
	logger.L.Info("InitPool", zap.Any("cfg", cfg))
	return pool.NewService(*cfg.PoolConfig)
}
