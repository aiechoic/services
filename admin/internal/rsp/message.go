package rsp

var Codes = map[int]string{
	CodeSuccess:           "success",
	CodeBadRequest:        "bad request",
	CodeUnauthorized:      "unauthorized",
	CodeServerError:       "server error",
	CodeVerifyCodeInvalid: "verify code invalid",
}

type Message struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"msg"`
	Data    any    `json:"data"`
}

const (
	CodeSuccess      = 200
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeServerError  = 500
)

const (
	CodeVerifyCodeInvalid = 1001
)

func Success(msg string, data any) *Message {
	return &Message{
		Code:    CodeSuccess,
		Type:    "success",
		Message: msg,
		Data:    data,
	}
}

func Error(code int, msg string) *Message {
	if msg == "" {
		msg = Codes[code]
	}
	if msg == "" {
		panic("unknown code")
	}
	return &Message{
		Code:    code,
		Type:    "error",
		Message: msg,
		Data:    nil,
	}
}

func Warning(msg string) *Message {
	return &Message{
		Code:    CodeSuccess,
		Type:    "warning",
		Message: msg,
		Data:    nil,
	}
}

func BadRequestError(err error) *Message {
	return Error(CodeBadRequest, err.Error())
}

func InternalServerError(err error) *Message {
	return Error(CodeServerError, err.Error())
}

func UnauthorizedWarning() *Message {
	return Warning("unauthorized")
}
