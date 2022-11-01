package workers

import (
	"fmt"
	"github.com/tony-zhuo/consume-fcm/internal/services"
)

type QueueWorker interface {
	NotificationConsume()
}

type queueWorker struct {
	queueService services.QueueService
}

func NewQueueService(queueService services.QueueService) QueueWorker {
	return &queueWorker{
		queueService: queueService,
	}
}

func (s *queueWorker) NotificationConsume() {
	fmt.Println("controller NotificationConsume start")
	s.queueService.Consumer()
}
