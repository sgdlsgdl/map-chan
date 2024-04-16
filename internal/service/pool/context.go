package pool

import (
	"context"
	"fmt"
	"github.com/sgdlsgdl/map-chan/pkg/baggage"
	"github.com/sgdlsgdl/map-chan/pkg/connector"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"go.uber.org/zap"
)

type WrappedConn struct {
	Conn   connector.Conn
	DestId int
	Addr   string
}

func (s *Service) getConnByContext(ctx context.Context, destId int) (connector.Conn, SendPacketCode, error) {
	dest, ok := s.discovery.GetDest(destId)
	if !ok {
		return nil, GetDestErr, fmt.Errorf("dest not found %d", destId)
	}

	wc, ok := baggage.GetVal(ctx).(*WrappedConn)
	if ok {
		if wc.DestId == destId && wc.Addr == dest.GetAddr() {
			if !wc.Conn.IsClosed() {
				return wc.Conn, Success, nil
			}
			logger.L.Warn("getConnByContext conn is closed", zap.Int("destId", destId), zap.String("addr", dest.GetAddr()),
				zap.String("conn", wc.Conn.Info()))
		} else {
			err := wc.Conn.Close()
			logger.L.Info("addr changed, close conn", zap.Int("destId", destId), zap.String("addr", dest.GetAddr()),
				zap.Int("oldDestId", wc.DestId), zap.String("oldAddr", wc.Addr), zap.String("conn", wc.Conn.Info()), zap.Error(err))
		}
	}

	p, err := s.connectorHandler.GetOrInitPool(dest)
	if err != nil {
		return nil, GetPoolErr, fmt.Errorf("GetOrInitPool %w", err)
	}
	conn, err := p.Get()
	if err != nil {
		return nil, GetConnErr, fmt.Errorf("get conn %w", err)
	}

	logger.L.Debug("set conn", zap.Int("destId", destId), zap.String("conn", conn.Info()))
	baggage.SetVal(ctx, &WrappedConn{
		Conn:   conn,
		DestId: destId,
		Addr:   dest.GetAddr(),
	})
	return conn, Success, nil
}

func (s *Service) destroyConnByContext(ctx context.Context, conn connector.Conn) {
	logger.L.Info("destroyConnByContext", zap.String("conn", conn.Info()))
	baggage.SetVal(ctx, nil)
	conn.Destroy()
}

func resetConn(ctx context.Context) {
	wc, ok := baggage.GetVal(ctx).(*WrappedConn)
	if ok {
		logger.L.Debug("reset conn", zap.String("conn", wc.Conn.Info()))
		baggage.SetVal(ctx, nil)
		err := wc.Conn.Close()
		if err != nil {
			logger.L.Error("close conn", zap.String("conn", wc.Conn.Info()), zap.Error(err))
		}
	}
}
