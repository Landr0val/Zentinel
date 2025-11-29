package handler

import (
	"net/http"
	"strconv"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/infrastructure/driving/http/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientHandler struct {
	clientUseCase port.ClientUseCase
}

func NewClientHandler(clientUseCase port.ClientUseCase) *ClientHandler {
	return &ClientHandler{
		clientUseCase: clientUseCase,
	}
}

func (h *ClientHandler) Create(c *gin.Context) {
	var req dto.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_REQUEST", "Invalid request body", err.Error()))
		return
	}

	client, err := h.clientUseCase.CreateClient(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to create client", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(client))
}

func (h *ClientHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid client ID", err.Error()))
		return
	}

	client, err := h.clientUseCase.GetClient(c.Request.Context(), id)
	if err != nil {
		// TODO: Handle not found specifically if use case returns specific error
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to get client", err.Error()))
		return
	}

	if client == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("NOT_FOUND", "Client not found", ""))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(client))
}

func (h *ClientHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid client ID", err.Error()))
		return
	}

	var req dto.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_REQUEST", "Invalid request body", err.Error()))
		return
	}

	client, err := h.clientUseCase.UpdateClient(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to update client", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(client))
}

func (h *ClientHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	clients, total, err := h.clientUseCase.ListClients(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to list clients", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(clients, page, pageSize, total))
}

func (h *ClientHandler) GetAccounts(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid client ID", err.Error()))
		return
	}

	accounts, err := h.clientUseCase.GetClientAccounts(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to get client accounts", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(accounts))
}
