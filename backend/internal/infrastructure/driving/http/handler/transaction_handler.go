package handler

import (
	"net/http"
	"strconv"
	"time"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"
	"zentinel/internal/infrastructure/driving/http/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionUseCase port.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase port.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_REQUEST", "Invalid request body", err.Error()))
		return
	}

	transaction, err := h.transactionUseCase.CreateTransaction(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to create transaction", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(transaction))
}

func (h *TransactionHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid transaction ID", err.Error()))
		return
	}

	transaction, err := h.transactionUseCase.GetTransaction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to get transaction", err.Error()))
		return
	}

	if transaction == nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("NOT_FOUND", "Transaction not found", ""))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(transaction))
}

func (h *TransactionHandler) List(c *gin.Context) {
	filter := repository.TransactionFilter{
		Page:     1,
		PageSize: 10,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			filter.PageSize = pageSize
		}
	}

	if fromDateStr := c.Query("from_date"); fromDateStr != "" {
		if fromDate, err := time.Parse(time.RFC3339, fromDateStr); err == nil {
			filter.FromDate = &fromDate
		}
	}

	if toDateStr := c.Query("to_date"); toDateStr != "" {
		if toDate, err := time.Parse(time.RFC3339, toDateStr); err == nil {
			filter.ToDate = &toDate
		}
	}

	if opTypeStr := c.Query("operation_type"); opTypeStr != "" {
		opType := enums.OperationType(opTypeStr)
		filter.OperationType = &opType
	}

	if channelStr := c.Query("channel"); channelStr != "" {
		channel := enums.Channel(channelStr)
		filter.Channel = &channel
	}

	if isFlaggedStr := c.Query("is_flagged"); isFlaggedStr != "" {
		if isFlagged, err := strconv.ParseBool(isFlaggedStr); err == nil {
			filter.IsFlagged = &isFlagged
		}
	}

	if accountIDStr := c.Query("account_id"); accountIDStr != "" {
		if accountID, err := uuid.Parse(accountIDStr); err == nil {
			filter.AccountID = &accountID
		}
	}

	transactions, total, err := h.transactionUseCase.ListTransactions(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to list transactions", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(transactions, filter.Page, filter.PageSize, total))
}

func (h *TransactionHandler) Analyze(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid transaction ID", err.Error()))
		return
	}

	transaction, err := h.transactionUseCase.AnalyzeTransaction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "Failed to analyze transaction", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(transaction))
}
