package worker_pool

import (
	"context"
	"log/slog"
	"sync"

	"github.com/zagvozdeen/ola/internal/logger"
)

type Task func() error

type WorkerPool struct {
	log *logger.Logger
	n   int
	ch  chan Task
}

func New(log *logger.Logger, n int, capacity int) *WorkerPool {
	return &WorkerPool{
		log: log,
		n:   n,
		ch:  make(chan Task, capacity),
	}
}

func (p *WorkerPool) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}

	for i := range p.n {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-p.ch:
					if !ok {
						return
					}
					if err := task(); err != nil {
						p.log.Error("Worker pool error", err, slog.Int("worker", i))
					}
				}
			}
		})
	}

	wg.Wait()
}

func (p *WorkerPool) Submit(fn Task) {
	p.ch <- fn
}
