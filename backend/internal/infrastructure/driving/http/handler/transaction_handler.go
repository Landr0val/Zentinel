package handler

import (
	"net/http"
	"strconv"
	"time"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/domain"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"
	"zentinel/internal/infrastructure/driving/http/dto/response"
	"zentinel/internal/infrastructure/driving/http/errorhandler"

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

// Create godoc
// @Summary      Create a new transaction
// @Description  Create a new transaction and analyze it for fraud
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateTransactionRequest true "Create Transaction Request"
// @Success      201  {object}  response.TransactionResponseWrapper
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/transactions [post]
// @Security     BearerAuth
func (h *TransactionHandler) Create(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleRequestError(c, err)
		return
	}

	transaction, err := h.transactionUseCase.CreateTransaction(c.Request.Context(), req)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(transaction))
}

// Get godoc
// @Summary      Get transaction by ID
// @Description  Get transaction details by ID
// @Tags         transactions
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  response.TransactionResponseWrapper
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/transactions/{id} [get]
// @Security     BearerAuth
func (h *TransactionHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid transaction ID", err.Error()))
		return
	}

	transaction, err := h.transactionUseCase.GetTransaction(c.Request.Context(), id)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	if transaction == nil {
		errorhandler.HandleDomainError(c, domain.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(transaction))
}

// List godoc
// @Summary      List transactions
// @Description  List transactions with filtering and pagination
// @Tags         transactions
// @Produce      json
// @Param        page            query     int     false  "Page number"
// @Param        page_size       query     int     false  "Page size"
// @Param        from_date       query     string  false  "From Date (RFC3339)"
// @Param        to_date         query     string  false  "To Date (RFC3339)"
// @Param        operation_type  query     string  false  "Operation Type"
// @Param        channel         query     string  false  "Channel"
// @Param        is_flagged      query     bool    false  "Is Flagged"
// @Param        account_id      query     string  false  "Account ID"
// @Success      200             {object}  response.TransactionListResponseWrapper
// @Failure      500             {object}  response.ErrorResponse
// @Router       /api/v1/transactions [get]
// @Security     BearerAuth
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
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(transactions, filter.Page, filter.PageSize, total))
}

// Analyze godoc
// @Summary      Analyze transaction
// @Description  Manually trigger fraud analysis for a transaction
// @Tags         transactions
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  response.TransactionResponseWrapper
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/transactions/{id}/analyze [post]
// @Security     BearerAuth
func (h *TransactionHandler) Analyze(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid transaction ID", err.Error()))
		return
	}

	transaction, err := h.transactionUseCase.AnalyzeTransaction(c.Request.Context(), id)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(transaction))
}
