package helper

import "github.com/Arshia-Izadyar/Fast-Gopher/src/api/validator"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

func GenerateResponse(data interface{}, success bool) *Response {
	return &Response{
		Success: success,
		Data:    data,
		Error:   "",
	}
}

func GenerateResponseWithError(err error, success bool) *Response {
	return &Response{
		Success: success,
		Data:    nil,
		Error:   err.Error(),
	}
}

func GenerateResponseWithValidationError(err error, success bool) *Response {
	ve := validator.GetValidationError(err)
	return &Response{
		Success: success,
		Data:    ve,
		Error:   err.Error(),
	}
}
