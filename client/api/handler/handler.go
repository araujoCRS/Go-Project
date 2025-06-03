package handler

import (
	"client/internal/service"
	shared "client/shared/models"
	"log"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerClient struct {
	service  service.ClientService
	rebbitMQ service.RabbitMQChannel
}

func NewHandlerClient(service *service.ClientService, rebbitMQ *service.RabbitMQChannel) *HandlerClient {
	return &HandlerClient{*service, *rebbitMQ}
}

// Create godoc
// @Summary Cria um novo cliente
// @Description Cria um novo cliente com os dados fornecidos.
// @Tags clients
// @Accept  json
// @Produce  json
// @Param   client body shared.Client true "Dados do Cliente"
// @Success 201 {object} shared.Client "Cliente criado com sucesso"
// @Failure 400 {object} map[string]string "Dados inválidos"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /client [post, put]
func (h *HandlerClient) Publish(ctx *gin.Context) {
	var client shared.Client
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := client.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.rebbitMQ.Publish(client)

	if err != nil {
		log.Printf("Erro ao publicar mensagem no RabbitMQ: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Contacte o suporte técnico"})
		return
	}

	ctx.JSON(http.StatusCreated, "Dados enviados para processamento")
}

// Get godoc
// @Summary Obtém um cliente pelo ID
// @Description Obtém os detalhes de um cliente específico pelo ID.
// @Tags clients
// @Accept  json
// @Produce  json
// @Param id path int true "ID do Cliente"
// @Success 200 {object} shared.Client "Cliente encontrado"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Cliente não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /client/{id} [get]
func (h *HandlerClient) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	client, err := h.service.Get(ctx, id)

	if err != nil {
		log.Printf("Erro: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Contacte o suporte técnico"})
		return
	}

	ctx.JSON(http.StatusOK, client)
}

// Delete godoc
// @Summary Deleta um cliente pelo ID
// @Description Deleta um cliente específico pelo ID.
// @Tags clients
// @Accept  json
// @Produce  json
// @Param id path int true "ID do Cliente"
// @Success 200 {object} map[string]string "Cliente deletado com sucesso"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Cliente não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /client/{id} [delete]
func (h *HandlerClient) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	deleted, err := h.service.Delete(ctx, id)

	if err != nil {
		log.Printf("Erro: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Contacte o suporte técnico"})
		return
	}

	if !deleted {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cliente deletado com sucesso"})
}
