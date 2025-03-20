package models

import (
	"net/http"
)

// ResponseParams struct
type ResponseParams struct {
	StatusCode     int
	SuccessMessage string
	ErrorMessage   string
	Content        any
}

// Response struct
type Response struct {
	Code           int    `json:"code"`
	SuccessMessage string `json:"success_message"`
	ErrorMessage   string `json:"error_message"`
	Content        any    `json:"content"`
}

// ResponseV2 struct for response
type ResponseV2 struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ResponseWithoutData struct
type ResponseWithoutData struct {
	Code           int    `json:"code"`
	SuccessMessage string `json:"successMessage"`
	ErrorMessage   string `json:"errorMessage"`
}

// NewResponse function for generate Response
func NewResponse(params ResponseParams) any {
	var response any
	var successMessage string

	if params.StatusCode >= 200 && params.StatusCode <= 299 {
		successMessage = "success"
	}

	if params.Content != nil {
		response = &Response{
			Code:           params.StatusCode,
			SuccessMessage: successMessage,
			ErrorMessage:   params.ErrorMessage,
			Content:        params.Content,
		}
	} else {
		response = &Response{
			Code:           params.StatusCode,
			SuccessMessage: successMessage,
			ErrorMessage:   params.ErrorMessage,
		}
	}

	return response
}

// ResponseSuccessWithData initiate Response
func ResponseSuccessWithData(data any) Response {
	response := Response{
		Code:           http.StatusOK,
		SuccessMessage: "success",
		Content:        data,
	}
	return response
}

// ResponseSuccess function for generate success response
func ResponseSuccess() Response {
	response := Response{
		Code:           http.StatusOK,
		SuccessMessage: "success",
	}
	return response
}

// ResponseError function for generate response error
func ResponseError(code int, msg string) Response {
	response := Response{
		Code:         code,
		ErrorMessage: msg,
	}
	return response
}
