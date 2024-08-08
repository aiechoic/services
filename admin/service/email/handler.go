package email

import (
	"github.com/aiechoic/services/admin/internal/rsp"
	"github.com/aiechoic/services/email/verify"
	"github.com/aiechoic/services/gins"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	s *verify.Generator
}

func NewHandlers(s *verify.Generator) *Handlers {
	return &Handlers{s: s}
}

func (s *Handlers) SendVerifyCode() gins.Handler {
	type SendVerifyCodeRequest struct {
		Email string `json:"email" binding:"required,email"`
	}
	type SendVerifyCodeResponse struct {
		TimeWait int64 `json:"time_wait"`
	}
	return gins.Handler{
		Request: gins.Request{
			Json: &SendVerifyCodeRequest{},
		},
		Response: gins.Response{
			Json:        &SendVerifyCodeResponse{},
			Description: "Return time(seconds) to wait for next send code",
		},
		Handler: func(c *gin.Context) {
			var req SendVerifyCodeRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			wait := s.s.GetWaitTime(req.Email)
			if wait > 0 {
				c.JSON(200, rsp.Success("ok", &SendVerifyCodeResponse{TimeWait: int64(wait.Seconds()) + 1}))
				return
			} else {
				_, err := s.s.GenerateCode(req.Email)
				if err != nil {
					c.JSON(200, rsp.InternalServerError(err))
				} else {
					wait = s.s.GetWaitTime(req.Email)
					c.JSON(200, rsp.Success("ok", &SendVerifyCodeResponse{TimeWait: int64(wait.Seconds()) + 1}))
				}
			}
		},
	}
}
