package {{.PackageName}}

import (
	"github.com/aiechoic/services/admin/internal/model"
	"gorm.io/gorm"
)

type {{.ModelName}} struct {
	ID int `json:"id" gorm:"primaryKey"`
}

type DB struct {
	*model.BaseDB[{{.ModelName}}]
}

func NewDB(gormDB *gorm.DB) (*DB, error) {
	baseDB, err := model.NewBaseDB[{{.ModelName}}](gormDB)
	if err != nil {
		return nil, err
	}
	return &DB{
		BaseDB: baseDB,
	}, nil
}