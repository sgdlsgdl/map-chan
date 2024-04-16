package pool

import (
	"context"
	"fmt"
	"github.com/sgdlsgdl/map-chan/pkg/define"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"github.com/sgdlsgdl/map-chan/pkg/map_chan"
	"go.uber.org/zap"
)

func (s *Service) Push(p define.Packet) error {
	if p == nil {
		return fmt.Errorf("packet invalid")
	}
	po := s.lazyInit(p.GetDestId())
	if po == nil {
		return fmt.Errorf("dest not found")
	}
	err := po.Send(p, p.GetKey())
	if err != nil {
		logger.L.Info("Push", zap.String("packet", p.Verbose()), zap.Error(err))
	}
	return err
}

func (s *Service) lazyInit(destId int) *map_chan.MapChan[define.Packet] {
	s.mcMapMutex.RLock()
	mc, ok := s.mcMap[destId]
	s.mcMapMutex.RUnlock()

	if ok {
		return mc
	}

	dest, ok := s.discovery.GetDest(destId)
	logger.L.Debug("lazyInit", zap.Int("destId", destId), zap.Any("dest", dest), zap.Bool("exist", ok))
	if !ok {
		return nil
	}

	s.mcMapMutex.Lock()
	defer s.mcMapMutex.Unlock()
	mc, ok = s.mcMap[destId]
	if ok {
		return mc
	}

	mc = map_chan.NewMapChan[define.Packet](fmt.Sprintf("%s-%d", dest.GetName(), destId),
		dest.GetConnNum(), s.chanSize, s.sendChanTimeout, func(ctx context.Context, p define.Packet) {
			err := s.writePacket(ctx, destId, p)
			if err != nil {
				logger.L.Error("writePacket", zap.String("packet", p.Verbose()), zap.Error(err))
			}
		}, resetConn)
	s.mcMap[destId] = mc
	return mc
}

func (s *Service) writePacket(ctx context.Context, destId int, p define.Packet) error {
	var (
		code = Success
		err  error
	)
	defer func() {
		if code == Success {
			logger.L.Debug("writePacket", zap.String("packet", p.Info()))
		} else {
			logger.L.Error("writePacket", zap.String("packet", p.Verbose()), zap.String("code", code.String()), zap.Error(err))
		}
	}()

	conn, code, err := s.getConnByContext(ctx, destId)
	if err != nil {
		return err
	}

	if err = conn.Write(p); err != nil {
		code = WritePacketErr
		s.destroyConnByContext(ctx, conn)
		return fmt.Errorf("writePacket %w", err)
	}

	return nil
}
