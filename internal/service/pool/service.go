package pool

import (
	"github.com/sgdlsgdl/map-chan/config/pool"
	"github.com/sgdlsgdl/map-chan/pkg/connector"
	"github.com/sgdlsgdl/map-chan/pkg/define"
	"github.com/sgdlsgdl/map-chan/pkg/discovery"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"github.com/sgdlsgdl/map-chan/pkg/map_chan"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Service struct {
	sendChanTimeout  time.Duration
	chanSize         int
	discovery        discovery.Discovery
	discoveryNotify  chan discovery.Discovery
	connectorHandler connector.Handler

	mcMap      map[int]*map_chan.MapChan[define.Packet]
	mcMapMutex sync.RWMutex
}

func NewService(cfg poolCfg.PoolConfig) *Service {
	s := &Service{
		sendChanTimeout:  cfg.SendChanTimeout,
		chanSize:         cfg.ChanSize,
		discovery:        cfg.Discovery,
		connectorHandler: cfg.ConnectorHandler,
		discoveryNotify:  make(chan discovery.Discovery),
		mcMap:            make(map[int]*map_chan.MapChan[define.Packet]),
		mcMapMutex:       sync.RWMutex{},
	}
	go s.handleDiscoveryChange()
	return s
}

func (s *Service) Reload(discovery discovery.Discovery) {
	select {
	case s.discoveryNotify <- discovery:
	default:
		logger.L.Error("Reload chan full")
	}
}

func (s *Service) handleDiscoveryChange() {
	for {
		select {
		case newD, ok := <-s.discoveryNotify:
			if !ok {
				logger.L.Info("handleDiscoveryChange chan closed")
				return
			}
			start := time.Now()
			changeMap := s.discovery.UpdateAndDiff(newD)
			s.handleDel(changeMap[discovery.ChangeEventDel])
			s.handleAddr(changeMap[discovery.ChangeEventAddr])
			s.handleNum(changeMap[discovery.ChangeEventNum])
			logger.L.Info("handleDiscoveryChange finished", zap.Duration("cost", time.Since(start)))
		}
	}
}

func (s *Service) handleDel(ls []discovery.Dest) {
	start := time.Now()
	defer func() {
		logger.L.Info("handleDel", zap.Any("ls", ls), zap.Duration("cost", time.Since(start)))
	}()
	if len(ls) == 0 {
		return
	}

	s.mcMapMutex.Lock()
	for _, d := range ls {
		mc, ok := s.mcMap[d.GetDestId()]
		if ok {
			mc.Close()
		}
	}
	s.mcMapMutex.Unlock()
	logger.L.Info("handleDel mcMap", zap.Any("ls", ls), zap.Duration("cost", time.Since(start)))

	s.connectorHandler.BatchClose(ls)
	logger.L.Info("handleDel connectorHandler", zap.Any("ls", ls), zap.Duration("cost", time.Since(start)))
}

func (s *Service) handleAddr(ls []discovery.Dest) {
	start := time.Now()
	defer func() {
		logger.L.Info("handleAddr", zap.Any("ls", ls), zap.Duration("cost", time.Since(start)))
	}()

	if len(ls) == 0 {
		return
	}
	// no need to reset
}

func (s *Service) handleNum(ls []discovery.Dest) {
	start := time.Now()
	defer func() {
		logger.L.Info("handleNum", zap.Any("ls", ls), zap.Duration("cost", time.Since(start)))
	}()

	if len(ls) == 0 {
		return
	}
	s.mcMapMutex.Lock()
	for _, d := range ls {
		mc, ok := s.mcMap[d.GetDestId()]
		if ok {
			mc.ChangeChanList(d.GetConnNum(), s.chanSize)
		}
	}
	s.mcMapMutex.Unlock()
}
