package user

import (
	"fmt"
	"github.com/aiechoic/services/gins"
)

func NewService(secret string) *gins.Service {
	db := NewDB()
	auth := NewJWTAuth(secret)
	h := NewHandlers(db, auth)
	return &gins.Service{
		Tag:         "user",
		Description: "User service",
		Path:        "/",
		Routes: []gins.Route{
			{
				Method:      "POST",
				Path:        "/login",
				Summary:     "Login",
				Description: fmt.Sprintf("default admin: name=\"%s\", password=\"%s\"", Admin.Name, Admin.Password),
				Handler:     h.Login(),
			},
			{
				Method:   "GET",
				Path:     "/user",
				Summary:  "Get user",
				Security: auth,
				Handler:  h.Get(),
			},
			{
				Method:   "POST",
				Path:     "/user",
				Summary:  "Create user",
				Security: auth,
				Handler:  h.Create(),
			},
			{
				Method:   "GET",
				Path:     "/users",
				Summary:  "List users",
				Security: auth,
				Handler:  h.List(),
			},
			{
				Method:   "PUT",
				Path:     "/user",
				Summary:  "Update user",
				Security: auth,
				Handler:  h.Update(),
			},
			{
				Method:   "DELETE",
				Path:     "/user",
				Summary:  "Delete user",
				Security: auth,
				Handler:  h.Delete(),
			},
		},
	}
}
