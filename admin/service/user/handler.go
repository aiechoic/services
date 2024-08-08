package user

import (
	"github.com/aiechoic/services/admin/internal/auth"
	"github.com/aiechoic/services/admin/internal/rsp"
	"github.com/aiechoic/services/email/verify"
	"github.com/aiechoic/services/gins"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	es   *verify.Generator
	auth *auth.Auth[User]
	db   *DB
}

func NewHandlers(es *verify.Generator, udb *DB, auth *auth.Auth[User]) *Handlers {
	return &Handlers{
		es:   es,
		db:   udb,
		auth: auth,
	}
}

func (a *Handlers) Register() gins.Handler {
	type RegisterRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Code     string `json:"code" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	return gins.Handler{
		Request: gins.Request{Json: RegisterRequest{}},
		Handler: func(c *gin.Context) {
			var req RegisterRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			valid, err := a.es.VerifyCode(req.Email, req.Code)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			if !valid {
				c.JSON(200, rsp.Error(rsp.CodeVerifyCodeInvalid, ""))
				return
			}
			user := &User{
				Email:    req.Email,
				Password: req.Password,
			}
			err = a.db.Create(user)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", user))
		},
	}
}

func (a *Handlers) Login() gins.Handler {
	type LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	return gins.Handler{
		Request: gins.Request{Json: LoginRequest{}},
		Response: gins.Response{
			Json: User{},
		},
		Handler: func(c *gin.Context) {
			var req LoginRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			user, err := a.db.Login(req.Email, req.Password)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", user))
		},
	}
}

func (a *Handlers) Logout() gins.Handler {
	return gins.Handler{
		Handler: func(c *gin.Context) {
			token := a.auth.GetToken(c)
			err := a.db.Logout(token)
			if err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", nil))
		},
	}
}

func (a *Handlers) ResetPassword() gins.Handler {
	type ResetPasswordRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Code     string `json:"code" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	return gins.Handler{
		Request: gins.Request{Json: ResetPasswordRequest{}},
		Handler: func(c *gin.Context) {
			var req ResetPasswordRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			valid, err := a.es.VerifyCode(req.Email, req.Code)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			if !valid {
				c.JSON(200, rsp.Error(rsp.CodeVerifyCodeInvalid, ""))
				return
			}
			err = a.db.ResetPassword(req.Email, req.Password)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", nil))
		},
	}
}

func (a *Handlers) GetProfile() gins.Handler {
	return gins.Handler{
		Response: gins.Response{Json: User{}},
		Handler: func(c *gin.Context) {
			user := a.auth.GetSession(c)
			if user != nil {
				c.JSON(200, rsp.Success("ok", user))
			}
		},
	}
}
