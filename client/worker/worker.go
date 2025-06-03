package worker

import (
	"client/internal/service"
	"log"
)

type Worker interface {
	Run() error
	Handler(msg []byte) error
}

type BaseWorker struct {
	rabbit *service.RabbitMQChannel
	//handler func (msg []byte) error
}

func (w *BaseWorker) Run(field Worker) error {
	deliveries, err := w.rabbit.Consume()

	if err != nil {
		log.Printf("Erro ao consumir mensagem: %v", err)
		return err
	}

	go func() {
		for msg := range deliveries {
			if err := field.Handler(msg.Body); err == nil {
				_ = msg.Ack(false) //successfully
			} else {
				log.Printf("Erro ao processar mensagem: %v", err)
				msg.Nack(false, true) // requeue
			}
		}
	}()

	return nil
}

func (w *BaseWorker) Handler(msg []byte) error {
	return nil
}
