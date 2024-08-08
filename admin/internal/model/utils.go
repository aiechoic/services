package model

import "gorm.io/gorm"

func GetTableName[T any](db *gorm.DB) (string, error) {
	var entity T
	stmt := &gorm.Statement{DB: db}
	err := stmt.Parse(&entity)
	if err != nil {
		return "", err
	}
	return stmt.Table, nil
}
