package main

import (
	"client/api/docs"
	"client/api/handler"
	"client/configs"
	"client/internal/service"
	"client/internal/service/repository"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// gracefulShutdown encerra o servidor HTTP e a aplicação de forma segura.
func GracefulShutdown(server *http.Server, channel chan os.Signal) {
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	<-channel

	log.Println("Server shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server shutdown successfully")
	os.Exit(0)
}

// @title Client API
// @version 1.0
// @description Uma API simples para cadastro.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	shutdownChan := make(chan os.Signal, 1)
	router := gin.Default()
	config, _ := configs.LoadConfig()
	strConn := config.Database.GetConnectionString()
	ctx := context.Background()

	bdContext := repository.NewDbContext(strConn, &ctx)
	repo := repository.NewClientRepository(bdContext.GetConnection())
	clientservice := service.NewClientService(repo)

	rabbitChannel, rabbitConnection, err := service.NewRabbitMQService(config.RabbitMQ, config.RabbitMQ.Queues["clients"])
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}

	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api")
	{
		clientHandler := handler.NewHandlerClient(&clientservice, rabbitChannel)
		v1.POST("/client", clientHandler.Publish)
		v1.PUT("/client", clientHandler.Publish)
		v1.GET("/client/:id", clientHandler.Get)
		v1.DELETE("/client/:id", clientHandler.Delete)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	rabbitChannel.ListenClose(rabbitConnection, func(err error) {
		if err != nil {
			log.Printf("RabbitMQ connection closed with error: %v", err)
		}
		shutdownChan <- syscall.SIGTERM
	})

	go func() {
		log.Printf("Servidor iniciado na porta %s", port)
		log.Printf("Acesse a documentação Swagger em http://localhost:%s/swagger/index.html", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	GracefulShutdown(server, shutdownChan)
}
