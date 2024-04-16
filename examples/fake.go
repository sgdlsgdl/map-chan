//go:build test
// +build test

package main

import (
	"fmt"
	"github.com/sgdlsgdl/map-chan/pkg/connector"
	"github.com/sgdlsgdl/map-chan/pkg/define"
	"github.com/sgdlsgdl/map-chan/pkg/discovery"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"go.uber.org/zap"
)

type FakeHandler struct {
}

func (f FakeHandler) GetOrInitPool(dest discovery.Dest) (connector.Pool, error) {
	return FakePool{dest}, nil
}

func (f FakeHandler) BatchClose([]discovery.Dest) {
	return
}

type FakePool struct {
	discovery.Dest
}

func (f FakePool) Get() (connector.Conn, error) {
	return FakeConn{f}, nil
}

func (f FakePool) Close() {
	return
}

func (f FakePool) Len() int {
	return 1
}

type FakeConn struct {
	discovery.Dest
}

func (f FakeConn) Write(p define.Packet) error {
	logger.L.Info("write packet", zap.Any("raw", p.GetRawPacket()), zap.Any("packet", p.Info()))
	return nil
}

func (f FakeConn) IsClosed() bool {
	return false
}

func (f FakeConn) Close() error {
	return nil
}

func (f FakeConn) Destroy() {
	return
}

func (f FakeConn) Info() string {
	return fmt.Sprintf("%d %s %s", f.GetDestId(), f.GetName(), f.GetAddr())
}
