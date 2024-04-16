//go:build test
// +build test

package main

import (
	"context"
	"fmt"
	concurrencyCfg "github.com/sgdlsgdl/map-chan/config/concurrency"
	poolCfg "github.com/sgdlsgdl/map-chan/config/pool"
	"github.com/sgdlsgdl/map-chan/facade"
	"github.com/sgdlsgdl/map-chan/pkg/define"
	"github.com/sgdlsgdl/map-chan/pkg/discovery"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"go.uber.org/zap"
	"time"
)

func main() {
	po := facade.InitPool(
		poolCfg.WithPoolSendChanTimeout(time.Millisecond*30),
		poolCfg.WithPoolChanSize(10),
		poolCfg.WithDiscovery(&discovery.DestMap{
			DestIdToDest: map[int]discovery.Dest{
				1000: discovery.DestServer{
					DestId:  1000,
					Name:    "def name",
					Addr:    "def addr",
					ConnNum: 10,
				},
			},
		}),
		poolCfg.WithConnectorHandler(FakeHandler{}),
	)

	fn := func(ctx context.Context, p define.Packet) {
		err := po.Push(p)
		if err != nil {
			logger.L.Error("push failed", zap.Error(err))
		}
	}
	con := facade.InitConcurrency(
		concurrencyCfg.WithConcurrencyFn(fn),
		concurrencyCfg.WithConcurrencySendChanTimeout(time.Millisecond*30),
		concurrencyCfg.WithConcurrencyChanSize(10),
		concurrencyCfg.WithConcurrencyChanNum(10),
	)

	go func() {
		count := 0
		t := time.NewTicker(time.Second * 11)
		for {
			select {
			case <-t.C:
				newD := &discovery.DestMap{
					DestIdToDest: map[int]discovery.Dest{
						1000: discovery.DestServer{
							DestId:  1000,
							Name:    fmt.Sprintf("%d", count),
							Addr:    fmt.Sprintf("%d", count),
							ConnNum: count,
						},
					},
				}
				po.Reload(newD)
				count++
			}
		}
	}()

	for i := 0; i < 2; i++ {
		i := i
		go func() {
			t := time.NewTicker(time.Second * 10)
			for {
				select {
				case <-t.C:
					err := con.Exec(&define.WrappedPacket{
						Packet: fmt.Sprintf("msg from %d", i),
						Key:    fmt.Sprintf("%d", i),
						DestId: 1000,
					})
					if err != nil {
						logger.L.Error("exec failed", zap.Error(err))
					}
				}
			}
		}()
	}

	select {}
}
