package concurrency

import (
	"context"
	"fmt"
	"github.com/sgdlsgdl/map-chan/config/concurrency"
	"github.com/sgdlsgdl/map-chan/pkg/define"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"github.com/sgdlsgdl/map-chan/pkg/map_chan"
	"go.uber.org/zap"
	"sync"
)

type Service struct {
	cfg       concurrencyCfg.ConcurrencyConfig
	cfgMutex  sync.RWMutex
	cfgNotify chan concurrencyCfg.ConcurrencyConfig

	mc *map_chan.MapChan[define.Packet]
}

func NewService(cfg concurrencyCfg.ConcurrencyConfig, fn func(ctx context.Context, p define.Packet)) *Service {
	s := &Service{
		cfg:       cfg,
		cfgMutex:  sync.RWMutex{},
		cfgNotify: make(chan concurrencyCfg.ConcurrencyConfig),
		mc:        map_chan.NewMapChan("concurrency", cfg.ChanNum, cfg.ChanSize, cfg.SendChanTimeout, fn, nil),
	}
	go s.handleConfigChange()
	return s
}

func (s *Service) ChangeChan(chanSize, chanNum int) {
	cfg := concurrencyCfg.ConcurrencyConfig{
		ChanSize: chanSize,
		ChanNum:  chanNum,
	}
	s.changeByConfig(cfg)
}

func (s *Service) changeByConfig(cfg concurrencyCfg.ConcurrencyConfig) {
	select {
	case s.cfgNotify <- cfg:
	default:
		logger.L.Error("changeByConfig chan full")
	}
}

func (s *Service) handleConfigChange() {
	for {
		select {
		case newCfg, ok := <-s.cfgNotify:
			if !ok {
				logger.L.Info("handleConfigChange chan closed")
				return
			}
			if newCfg.ChanNum == s.cfg.ChanNum && newCfg.ChanSize == s.cfg.ChanSize {
				logger.L.Info("handleConfigChange no change", zap.Any("cfg", newCfg))
				continue
			}
			s.cfgMutex.Lock()
			s.cfg = newCfg
			s.cfgMutex.Unlock()
			s.mc.ChangeChanList(newCfg.ChanNum, newCfg.ChanSize)
			logger.L.Info("handleConfigChange finished", zap.Any("cfg", newCfg))
		}
	}
}

func (s *Service) Exec(p define.Packet) error {
	if p == nil {
		return fmt.Errorf("exec packet invalid")
	}
	err := s.mc.Send(p, p.GetKey())
	if err != nil {
		logger.L.Info("Exec", zap.String("packet", p.Verbose()), zap.Error(err))
	}
	return err
}
