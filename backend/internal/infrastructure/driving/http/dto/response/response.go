package response

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Data       []T            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

func NewSuccessResponse[T any](data T) SuccessResponse[T] {
	return SuccessResponse[T]{
		Data: data,
	}
}

func NewPaginatedResponse[T any](data []T, page, pageSize int, totalItems int64) PaginatedResponse[T] {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((totalItems + int64(pageSize) - 1) / int64(pageSize))
	}

	return PaginatedResponse[T]{
		Data: data,
		Pagination: PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}
}

func NewErrorResponse(code, message, details string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetails{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}
