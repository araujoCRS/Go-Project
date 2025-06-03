package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"client/configs"
	"client/internal/service"
	"client/internal/service/repository"
	"client/worker"
)

func main() {
	ctx := context.Background()
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	strConn := config.Database.GetConnectionString()

	bdContext := repository.NewDbContext(strConn, &ctx)
	repo := repository.NewClientRepository(bdContext.GetConnection())
	clientservice := service.NewClientService(repo)

	rabbitChannel, conn, err := service.NewRabbitMQService(config.RabbitMQ, config.RabbitMQ.Queues["clients"])
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}

	defer conn.Close()
	defer rabbitChannel.Close()

	var workerClient worker.Worker = worker.NewWorkerClient(&ctx, rabbitChannel, &clientservice)

	if err := workerClient.Run(); err != nil {
		log.Fatalf("Erro ao iniciar o worker: %v", err)
	}

	log.Println("Worker iniciado e aguardando mensagens...")

	// Aguarda sinal para encerrar
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Println("Worker finalizado.")
}
