package v1

type ResponseMessage struct {
	Code    ErrCode     `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
	Debug   interface{} `json:"debug,omitempty"`
}

func NewResp(code ErrCode, message string, payload ...interface{}) *ResponseMessage {
	if len(payload) > 0 {
		return &ResponseMessage{code, message, payload[0], nil}
	} else {
		return &ResponseMessage{code, message, nil, nil}
	}
}

func NewDebugResp(code ErrCode, message string, debug interface{}, payload ...interface{}) *ResponseMessage {
	if len(payload) > 0 {
		return &ResponseMessage{code, message, payload[0], debug}
	} else {
		return &ResponseMessage{code, message, nil, debug}
	}
}

type LoginResponse struct {
	Token string `json:"token"`
}
