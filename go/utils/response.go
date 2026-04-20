package utils

type Response[T any] struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func ErrorResponse(err error, message string, code int) map[string]any {
	return map[string]any{
		"code":    code,
		"error":   err.Error(),
		"message": message,
		"success": false,
	}
}

func SuccessResponse[T any](data T) Response[T] {
	return Response[T]{
		Code:    200,
		Message: "success",
		Data:    data,
		Success: true,
	}
}
