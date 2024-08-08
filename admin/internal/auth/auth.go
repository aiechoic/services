package auth

import (
	"context"
	"github.com/aiechoic/services/admin/internal/model"
	"github.com/aiechoic/services/admin/internal/rsp"
	"github.com/gin-gonic/gin"
)

type Auth[T any] struct {
	db             *model.AuthDB[T]
	securityScheme string
	headerKey      string
}

func NewAuth[T any](db *model.AuthDB[T], securitySchema, headerKey string) *Auth[T] {
	return &Auth[T]{
		db:             db,
		securityScheme: securitySchema,
		headerKey:      headerKey,
	}
}

func (m *Auth[T]) SecurityScheme() []map[string][]string {
	return []map[string][]string{
		{
			m.securityScheme: {},
		},
	}
}

func (m *Auth[T]) Auth(c *gin.Context) {
	admin := m.GetSession(c)
	if admin == nil {
		c.JSON(200, rsp.Error(rsp.CodeUnauthorized, ""))
		c.Abort()
		return
	}
	c.Next()
}

func (m *Auth[T]) GetToken(c *gin.Context) string {
	return c.GetHeader(m.headerKey)
}

func (m *Auth[T]) GetSession(c *gin.Context) *T {
	token := m.GetToken(c)
	if token == "" {
		return nil
	}
	if v, ok := c.Get(token); ok {
		if t, ok := v.(*T); ok {
			return t
		}
	}
	t, err := m.db.GetByToken(context.Background(), token)
	if err != nil {
		return nil
	}
	c.Set(token, t)
	return t
}
