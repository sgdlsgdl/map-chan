package map_chan

import (
	"container/list"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sgdlsgdl/map-chan/pkg/baggage"
	"github.com/sgdlsgdl/map-chan/pkg/logger"
	"go.uber.org/zap"
	"hash/fnv"
	"sync"
	"time"
)

const (
	DefaultChanNum         = 16
	DefaultChanSize        = 1000
	DefaultSendChanTimeout = 1 * time.Second
)

type MapChan[T any] struct {
	name            string
	hash            func(key string) int
	list            *list.List
	sendChanTimeout time.Duration
	chanNum         int
	chanSize        int
	mutex           sync.RWMutex
	fn              func(context.Context, T)
	cleanFn         func(context.Context)
	duration        time.Duration
}

type ChanElement[T any] struct {
	ctx       context.Context
	ch        chan T
	name      string
	uuid      string
	startTime int64
	quit      chan struct{}
	fn        func(context.Context, T)
	cleanFn   func(context.Context)
}

func NewMapChan[T any](name string, chanNum, chanSize int, sendChanTimeout time.Duration, fn func(ctx context.Context, t T), cleanFn func(context.Context)) *MapChan[T] {
	if sendChanTimeout <= 0 {
		sendChanTimeout = DefaultSendChanTimeout
	}
	mc := &MapChan[T]{
		name: name,
		hash: func(key string) int {
			h := fnv.New32a()
			h.Write([]byte(key))
			hashValue := h.Sum32()
			return int(hashValue)
		},
		list:            list.New(),
		sendChanTimeout: sendChanTimeout,
		chanNum:         chanNum,
		chanSize:        chanSize,
		mutex:           sync.RWMutex{},
		fn:              fn,
		cleanFn:         cleanFn,
	}
	mc.setChanList()
	return mc
}

func (c *MapChan[T]) setChanList() {
	currentLen := c.list.Len()
	for i := 0; i < c.chanNum; i++ {
		id, _ := uuid.NewV4()
		ce := &ChanElement[T]{
			ctx:       baggage.NewBaggage(),
			ch:        make(chan T, c.chanSize),
			name:      c.name,
			uuid:      id.String(),
			startTime: time.Now().UnixNano(),
			quit:      make(chan struct{}),
			fn:        c.fn,
			cleanFn:   c.cleanFn,
		}
		go ce.startChanElement()
		c.list.PushBack(ce)
	}
	for i := 0; i < currentLen; i++ {
		e := c.list.Front()
		c.list.Remove(e)
		ce := e.Value.(*ChanElement[T])
		close(ce.quit)
	}
}

func (ce *ChanElement[T]) startChanElement() {
	logger.L.Debug("chan start", zap.String("name", ce.name), zap.String("uuid", ce.uuid), zap.Int64("startTime", ce.startTime))
	for {
		select {
		case <-ce.quit:
			logger.L.Debug("chan quit", zap.String("name", ce.name), zap.String("uuid", ce.uuid), zap.Int64("startTime", ce.startTime))
			for {
				select {
				case val := <-ce.ch:
					if ce.fn != nil {
						ce.fn(ce.ctx, val)
					}
				default:
					if ce.cleanFn != nil {
						ce.cleanFn(ce.ctx)
					}
					logger.L.Debug("chan quit finished", zap.String("name", ce.name), zap.String("uuid", ce.uuid), zap.Int64("startTime", ce.startTime))
					return
				}
			}
		case val := <-ce.ch:
			if ce.fn != nil {
				ce.fn(ce.ctx, val)
			}
		}
	}
}

func (c *MapChan[T]) Send(v T, key string) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	ce := c.findByIndex(c.hash(key))
	if ce == nil {
		return fmt.Errorf("no valid ChanElement")
	}
	select {
	case ce.ch <- v:
	case <-time.After(c.sendChanTimeout):
		return fmt.Errorf("chan block name %v uuid %v startTime %v val %v", ce.name, ce.uuid, ce.startTime, v)
	}
	return nil
}

func (c *MapChan[T]) findByIndex(index int) *ChanElement[T] {
	l := c.list.Len()
	if l == 0 {
		return nil
	}
	index = index % l
	e := c.list.Front()
	for i := 0; i < index; i++ {
		e = e.Next()
	}
	ch := e.Value.(*ChanElement[T])
	return ch
}

func (c *MapChan[T]) ChangeChanList(chanNum, chanSize int) {
	if chanNum == 0 {
		chanNum = DefaultChanNum
	}
	if chanSize == 0 {
		chanSize = DefaultChanSize
	}
	start := time.Now()
	defer func() {
		logger.L.Info("ChangeChanList", zap.String("name", c.name), zap.Duration("cost", time.Now().Sub(start)))
	}()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.chanNum == chanNum && c.chanSize == chanSize {
		return
	}
	c.chanNum, c.chanSize = chanNum, chanSize
	c.setChanList()
}

func (c *MapChan[T]) Close() {
	start := time.Now()
	defer func() {
		logger.L.Info("MapChan closed", zap.String("name", c.name), zap.Duration("cost", time.Now().Sub(start)))
	}()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i := 0; i < c.list.Len(); i++ {
		e := c.list.Front()
		c.list.Remove(e)
		ce := e.Value.(*ChanElement[T])
		close(ce.quit)
	}
}
