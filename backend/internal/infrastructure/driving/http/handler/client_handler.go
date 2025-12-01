package handler

import (
	"net/http"
	"strconv"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/domain"
	"zentinel/internal/infrastructure/driving/http/dto/response"
	"zentinel/internal/infrastructure/driving/http/errorhandler"

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

// Create godoc
// @Summary      Create a new client
// @Description  Create a new client with initial account
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateClientRequest true "Create Client Request"
// @Success      201  {object}  response.ClientResponseWrapper
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/clients [post]
// @Security     BearerAuth
func (h *ClientHandler) Create(c *gin.Context) {
	var req dto.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleRequestError(c, err)
		return
	}

	client, err := h.clientUseCase.CreateClient(c.Request.Context(), req)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(client))
}

// Get godoc
// @Summary      Get client by ID
// @Description  Get client details by ID
// @Tags         clients
// @Produce      json
// @Param        id   path      string  true  "Client ID"
// @Success      200  {object}  response.ClientResponseWrapper
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/clients/{id} [get]
// @Security     BearerAuth
func (h *ClientHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid client ID format", err.Error()))
		return
	}

	client, err := h.clientUseCase.GetClient(c.Request.Context(), id)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	if client == nil {
		errorhandler.HandleDomainError(c, domain.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(client))
}

// Update godo
// @Summary      Update client
// @Description  Update client details
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        id      path      string                   true  "Client ID"
// @Param        request body      dto.UpdateClientRequest  true  "Update Client Request"
// @Success      200     {object}  response.ClientResponseWrapper
// @Failure      400     {object}  response.ErrorResponse
// @Failure      404     {object}  response.ErrorResponse
// @Failure      500     {object}  response.ErrorResponse
// @Router       /api/v1/clients/{id} [put]
// @Security     BearerAuth
func (h *ClientHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid client ID format", err.Error()))
		return
	}

	var req dto.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleRequestError(c, err)
		return
	}

	client, err := h.clientUseCase.UpdateClient(c.Request.Context(), id, req)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(client))
}

// List godoc
// @Summary      List clients
// @Description  List clients with pagination
// @Tags         clients
// @Produce      json
// @Param        page       query     int  false  "Page number"
// @Param        page_size  query     int  false  "Page size"
// @Success      200        {object}  response.ClientListResponseWrapper
// @Failure      500        {object}  response.ErrorResponse
// @Router       /api/v1/clients [get]
// @Security     BearerAuth
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
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(clients, page, pageSize, total))
}

// GetAccounts godoc
// @Summary      Get client accounts
// @Description  Get all accounts associated with a client
// @Tags         clients
// @Produce      json
// @Param        id   path      string  true  "Client ID"
// @Success      200  {object}  response.AccountListWrapper
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/clients/{id}/accounts [get]
// @Security     BearerAuth
func (h *ClientHandler) GetAccounts(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid client ID format", err.Error()))
		return
	}

	accounts, err := h.clientUseCase.GetClientAccounts(c.Request.Context(), id)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(accounts))
}
