package event_bus

import (
	"context"
	"maps"
	"sync"

	"github.com/zagvozdeen/ola/internal/worker_pool"
)

type Handler[T any] func(context.Context, T) error

type Event[T any] struct {
	subs   map[uint64]Handler[T]
	pool   *worker_pool.WorkerPool
	nextID uint64
	//onError func(error)
	mu sync.RWMutex
}

func NewEvent[T any](pool *worker_pool.WorkerPool) *Event[T] {
	return &Event[T]{
		subs: map[uint64]Handler[T]{},
		pool: pool,
		//onError: onError,
	}
}

func (e *Event[T]) Subscribe(handler Handler[T]) (unsubscribe func()) {
	e.mu.Lock()
	defer e.mu.Unlock()

	id := e.nextID
	e.nextID++
	e.subs[id] = handler

	return func() {
		e.mu.Lock()
		defer e.mu.Unlock()
		delete(e.subs, id)
	}
}

func (e *Event[T]) Publish(ctx context.Context, event T) {
	subs := map[uint64]Handler[T]{}
	e.mu.RLock()
	maps.Copy(subs, e.subs)
	e.mu.RUnlock()

	for _, fn := range subs {
		e.pool.Submit(func() error {
			return fn(ctx, event)
		})
	}
}
