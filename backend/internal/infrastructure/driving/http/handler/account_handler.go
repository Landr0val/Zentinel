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

type AccountHandler struct {
	accountUseCase port.AccountUseCase
}

func NewAccountHandler(accountUseCase port.AccountUseCase) *AccountHandler {
	return &AccountHandler{
		accountUseCase: accountUseCase,
	}
}

// Create godoc
// @Summary      Create a new account
// @Description  Create a new account for a client
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateAccountRequest true "Create Account Request"
// @Success      201  {object}  response.AccountResponseWrapper
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/accounts [post]
// @Security     BearerAuth
func (h *AccountHandler) Create(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleRequestError(c, err)
		return
	}

	account, err := h.accountUseCase.CreateAccount(c.Request.Context(), req)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(account))
}

// Get godoc
// @Summary      Get account by ID or Number
// @Description  Get account details by ID or Account Number
// @Tags         accounts
// @Produce      json
// @Param        id   path      string  true  "Account ID or Number"
// @Success      200  {object}  response.AccountResponseWrapper
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/accounts/{id} [get]
// @Security     BearerAuth
func (h *AccountHandler) Get(c *gin.Context) {
	idStr := c.Param("id")

	var account *dto.AccountResponse
	var err error

	if id, parseErr := uuid.Parse(idStr); parseErr == nil {
		account, err = h.accountUseCase.GetAccount(c.Request.Context(), id)
	} else {
		account, err = h.accountUseCase.GetAccountByNumber(c.Request.Context(), idStr)
	}

	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	if account == nil {
		errorhandler.HandleDomainError(c, domain.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(account))
}

// Update godoc
// @Summary      Update account
// @Description  Update account details (e.g. status)
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id      path      string                  true  "Account ID"
// @Param        request body      dto.UpdateAccountRequest true  "Update Account Request"
// @Success      200     {object}  response.AccountResponseWrapper
// @Failure      400     {object}  response.ErrorResponse
// @Failure      404     {object}  response.ErrorResponse
// @Failure      500     {object}  response.ErrorResponse
// @Router       /api/v1/accounts/{id} [put]
// @Security     BearerAuth
func (h *AccountHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_ID", "Invalid account ID", err.Error()))
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleRequestError(c, err)
		return
	}

	account, err := h.accountUseCase.UpdateAccount(c.Request.Context(), id, req)
	if err != nil {
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(account))
}

// List godoc
// @Summary      List accounts
// @Description  List accounts with pagination
// @Tags         accounts
// @Produce      json
// @Param        page       query     int  false  "Page number"
// @Param        page_size  query     int  false  "Page size"
// @Success      200        {object}  response.AccountListResponseWrapper
// @Failure      500        {object}  response.ErrorResponse
// @Router       /api/v1/accounts [get]
// @Security     BearerAuth
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
		errorhandler.HandleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(accounts, page, pageSize, total))
}
