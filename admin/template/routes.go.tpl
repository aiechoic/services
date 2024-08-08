package {{.PackageName}}

import (
	"github.com/aiechoic/services/admin/service/admin"
	"github.com/aiechoic/services/database/gorm"
	"github.com/aiechoic/services/gins"
	"github.com/aiechoic/services/ioc"
)

func NewService(c *ioc.Container) *gins.Service {
	gormDB := gorm.GetGormDB(c)
	db, err := NewDB(gormDB)
	if err != nil {
		panic(err)
	}
	hs := NewHandlers(db)
	adminAuth := admin.NewAuth(c)
	return &gins.Service{
		Tag:         "{{.ModelName}} Service",
		Path:        "/{{.PackageName}}",
		Routes: []gins.Route{
			{
				Method:   "POST",
				Path:     "/get",
				Security: adminAuth,
				Handler:  hs.Get(),
			},
			{
				Method:   "POST",
				Path:     "/create",
				Security: adminAuth,
				Handler:  hs.Create(),
			},
			{
				Method:   "POST",
				Path:     "/full_update",
				Security: adminAuth,
				Handler:  hs.FullUpdate(),
			},
			{
				Method:   "POST",
				Path:     "/partial_update",
				Security: adminAuth,
				Handler:  hs.PartialUpdate(),
			},
			{
				Method:   "POST",
				Path:     "/delete",
				Security: adminAuth,
				Handler:  hs.Delete(),
			},
			{
				Method:   "POST",
				Path:     "/count",
				Security: adminAuth,
				Handler:  hs.Count(),
			},
			{
				Method:   "POST",
				Path:     "/list",
				Security: adminAuth,
				Handler:  hs.List(),
			},
		},
	}
}