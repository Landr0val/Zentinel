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

type AccountHandler struct {
	accountUseCase port.AccountUseCase
}

func NewAccountHandler(accountUseCase port.AccountUseCase) *AccountHandler {
	return &AccountHandler{
		accountUseCase: accountUseCase,
	}
}

func (h *AccountHandler) Create(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_REQUEST", "Invalid request body", err.Error()))
		return
	}

	account, err := h.accountUseCase.CreateAccount(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to create account", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(account))
}

func (h *AccountHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid account ID", err.Error()))
		return
	}

	account, err := h.accountUseCase.GetAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to get account", err.Error()))
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("NOT_FOUND", "Account not found", ""))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(account))
}

func (h *AccountHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid account ID", err.Error()))
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_REQUEST", "Invalid request body", err.Error()))
		return
	}

	account, err := h.accountUseCase.UpdateAccount(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to update account", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(account))
}

func (h *AccountHandler) List(c *gin.Context) {
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

	accounts, total, err := h.accountUseCase.ListAccounts(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to list accounts", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(accounts, page, pageSize, total))
}
