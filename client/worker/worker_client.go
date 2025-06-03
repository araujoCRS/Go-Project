package worker

import (
	"client/internal/service"
	shared "client/shared/models"
	"context"
	"encoding/json"
	"log"
)

type WorkerClient struct {
	*BaseWorker
	service service.ClientService
	ctx     context.Context
}

func (w *WorkerClient) Handler(msg []byte) error {
	var client shared.Client
	if err := json.Unmarshal(msg, &client); err != nil {
		log.Printf("Erro ao deserializar mensagem do cliente: %v", err)
		return err
	}
	log.Printf("Processando mensagem do cliente: %s", client.Cpf)
	if _, err := w.service.Save(w.ctx, client); err != nil {
		log.Printf("Erro ao salvar cliente: %v", err)
		return err
	}
	log.Printf("Cliente processado com sucesso: %s", client.Cpf)
	return nil
}

func (w *WorkerClient) Run() error {
	return w.BaseWorker.Run(w)
}

func NewWorkerClient(ctx *context.Context, rabbit *service.RabbitMQChannel, service *service.ClientService) *WorkerClient {
	base := &BaseWorker{rabbit: rabbit}
	return &WorkerClient{
		BaseWorker: base,
		service:    *service,
		ctx:        *ctx,
	}
}
