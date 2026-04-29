package dto

type Response struct {
	Message  string      `json:"message,omitempty"`
	Response int         `json:"response,omitempty"`
	Result   interface{} `json:"result,omitempty"`
}
