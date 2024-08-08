package email

import (
	"github.com/aiechoic/services/email/verify"
	"github.com/aiechoic/services/gins"
	"github.com/aiechoic/services/ioc"
)

func NewService(c *ioc.Container) *gins.Service {
	emailGenerator := verify.GetGenerator(c)
	sender := NewHandlers(emailGenerator)
	return &gins.Service{
		Tag:         "Email",
		Description: "Email service",
		Path:        "/email",
		Routes: []gins.Route{
			{
				Method:      "POST",
				Path:        "/send_verify_code",
				Description: "Send verification code to email, return time(seconds) to wait for next send code",
				Handler:     sender.SendVerifyCode(),
			},
		},
	}
}
