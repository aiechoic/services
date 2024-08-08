package user

import (
	"github.com/aiechoic/services/gins"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	db   *DB
	auth *JWTAuth
}

func NewHandlers(db *DB, auth *JWTAuth) *Handlers {
	return &Handlers{
		db:   db,
		auth: auth,
	}
}

func (h *Handlers) Get() gins.Handler {
	type getUserRequest struct {
		Id int `form:"id" binding:"required" description:"User ID"`
	}
	return gins.Handler{
		Request: gins.Request{
			Query: getUserRequest{},
		},
		Response: gins.Response{
			Json: User{},
		},
		Handler: func(c *gin.Context) {
			req := &getUserRequest{}
			if err := c.ShouldBindQuery(req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, h.db.GetUser(req.Id))
		},
	}
}

func (h *Handlers) List() gins.Handler {
	return gins.Handler{
		Response: gins.Response{
			Json: []*User{},
		},
		Handler: func(c *gin.Context) {
			c.JSON(200, h.db.GetUsers())
		},
	}
}

func (h *Handlers) Create() gins.Handler {
	return gins.Handler{
		Request: gins.Request{
			Form: User{},
		},
		Response: gins.Response{
			Json: User{},
		},
		Handler: func(c *gin.Context) {
			user := &User{}
			if err := c.ShouldBind(user); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			user, err := h.db.CreateUser(user)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, user)
		},
	}
}

func (h *Handlers) Update() gins.Handler {
	return gins.Handler{
		Request: gins.Request{
			Form: User{},
		},
		Response: gins.Response{
			Json: User{},
		},
		Handler: func(c *gin.Context) {
			user := &User{}
			if err := c.ShouldBind(user); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			user, err := h.db.UpdateUser(user)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, user)
		},
	}
}

func (h *Handlers) Delete() gins.Handler {
	type deleteUserRequest struct {
		Id int `form:"id" binding:"required" description:"User ID"`
	}
	return gins.Handler{
		Request: gins.Request{
			Query: deleteUserRequest{},
		},
		Response: gins.Response{
			Json: []*User{},
		},
		Handler: func(c *gin.Context) {
			req := &deleteUserRequest{}
			if err := c.ShouldBindQuery(req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			h.db.DeleteUser(req.Id)
			users := h.db.GetUsers()
			c.JSON(200, users)
		},
	}
}

func (h *Handlers) Login() gins.Handler {
	type loginRequest struct {
		Name     string `form:"name" binding:"required" description:"User name"`
		Password string `form:"password" binding:"required" description:"User password"`
	}
	type loginResponse struct {
		Token string `json:"token"`
	}
	return gins.Handler{
		Request: gins.Request{
			Form: loginRequest{},
		},
		Response: gins.Response{
			Json: loginResponse{},
		},
		Handler: func(c *gin.Context) {
			req := &loginRequest{}
			if err := c.ShouldBind(req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			user := h.db.Login(req.Name, req.Password)
			if user == nil {
				c.JSON(401, gin.H{"error": "Invalid name or password"})
				return
			}
			token, err := h.auth.GenerateToken(user)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"token": token})
		},
	}
}
