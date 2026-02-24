package event_bus

import (
	"github.com/zagvozdeen/ola/internal/store/models"
	"github.com/zagvozdeen/ola/internal/worker_pool"
)

type EventBus struct {
	OrderCreated    *Event[*models.Order]
	FeedbackCreated *Event[*models.Feedback]
}

func New(pool *worker_pool.WorkerPool) *EventBus {
	return &EventBus{
		OrderCreated:    NewEvent[*models.Order](pool),
		FeedbackCreated: NewEvent[*models.Feedback](pool),
	}
}
