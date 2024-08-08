package admin

import (
	"github.com/aiechoic/services/admin/internal/auth"
	"github.com/aiechoic/services/admin/internal/rsp"
	"github.com/aiechoic/services/admin/service/user"
	"github.com/aiechoic/services/email/verify"
	"github.com/aiechoic/services/gins"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	es   *verify.Generator
	db   *DB
	auth *auth.Auth[Admin]
	udb  *user.DB
}

func NewHandlers(es *verify.Generator, db *DB, udb *user.DB, auth *auth.Auth[Admin]) *Handlers {
	return &Handlers{
		es:   es,
		db:   db,
		udb:  udb,
		auth: auth,
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
			Json: Admin{},
		},
		Handler: func(c *gin.Context) {
			var req LoginRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			admin, err := a.db.Login(req.Email, req.Password)
			if err != nil {
				c.JSON(200, rsp.Warning(err.Error()))
				return
			}
			c.JSON(200, rsp.Success("ok", admin))
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
		Response: gins.Response{Json: Admin{}},
		Handler: func(c *gin.Context) {
			admin := a.auth.GetSession(c)
			if admin != nil {
				c.JSON(200, rsp.Success("ok", admin))
			}
		},
	}
}

func (a *Handlers) CountUser() gins.Handler {
	type CountUserResponse struct {
		Count int64 `json:"count"`
	}
	return gins.Handler{
		Response: gins.Response{Json: &CountUserResponse{}},
		Handler: func(c *gin.Context) {
			count, err := a.udb.Count()
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", CountUserResponse{Count: count}))
		},
	}
}

func (a *Handlers) FindUsers() gins.Handler {
	type ListUserRequest struct {
		By     string `json:"by" description:"id/email, default is id"`
		Desc   bool   `json:"desc" description:"false-asc true-desc"`
		Limit  int    `json:"limit" binding:"required"`
		Offset int    `json:"offset"`
	}
	return gins.Handler{
		Request: gins.Request{
			Json: &ListUserRequest{},
		},
		Response: gins.Response{
			Json: []*user.User{},
		},
		Handler: func(c *gin.Context) {
			var req ListUserRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			if req.By == "" {
				req.By = "id"
			}
			users, err := a.udb.FindByColumn(req.By, req.Desc, req.Limit, req.Offset)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", users))
		},
	}
}

func (a *Handlers) SetUserVipTime() gins.Handler {
	type SetUserVIPTimeRequest struct {
		Email string `json:"email" binding:"required,email"`
		Time  int64  `json:"time" binding:"required"`
	}
	return gins.Handler{
		Request:  gins.Request{Json: SetUserVIPTimeRequest{}},
		Response: gins.Response{Json: &user.User{}},
		Handler: func(c *gin.Context) {
			var req SetUserVIPTimeRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			u, err := a.udb.SetVipExpire(req.Email, req.Time)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", u))
		},
	}
}
