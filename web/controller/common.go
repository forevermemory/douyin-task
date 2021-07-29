package controller

// Response
type Response struct {
	Code int         `json:"code,omitempty"`
	Msg  int         `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
