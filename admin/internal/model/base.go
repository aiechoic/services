package model

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type BaseDB[T any] struct {
	db *gorm.DB
}

func NewBaseDB[T any](db *gorm.DB) (*BaseDB[T], error) {
	var entity T
	err := db.AutoMigrate(&entity)
	return &BaseDB[T]{
		db: db,
	}, err
}

func (b *BaseDB[T]) Create(entity *T) error {
	if err := b.db.Create(entity).Error; err != nil {
		return err
	}
	return nil
}

func (b *BaseDB[T]) GetByColumn(uniqueColumn string, value interface{}) (*T, error) {
	var entity T
	err := b.db.Where(uniqueColumn+" = ?", value).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (b *BaseDB[T]) FullUpdateByColumn(entity *T, uniqueColumn string, value interface{}) (*T, error) {
	var model T
	err := b.db.Model(&model).Where(uniqueColumn+" = ?", value).Select("*").Updates(entity).Error
	if err != nil {
		return nil, err
	}
	err = b.db.Where(uniqueColumn+" = ?", value).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (b *BaseDB[T]) PartialUpdateByColumn(entity *T, uniqueColumn string, value interface{}) (*T, error) {
	var model T
	err := b.db.Model(&model).Where(uniqueColumn+" = ?", value).Updates(entity).Error
	if err != nil {
		return nil, err
	}
	err = b.db.Where(uniqueColumn+" = ?", value).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (b *BaseDB[T]) DeleteByColumn(uniqueColumn string, value interface{}) (rowsAffected int64) {
	var entity T
	return b.db.Where(uniqueColumn+" = ?", value).Delete(&entity).RowsAffected
}

func (b *BaseDB[T]) FindByColumn(column string, desc bool, limit, offset int) ([]*T, error) {
	var entities []*T
	err := b.db.Order(clause.OrderByColumn{
		Column: clause.Column{Name: column},
		Desc:   desc,
	}).Limit(limit).Offset(offset).Find(&entities).Error
	return entities, err
}

func (b *BaseDB[T]) Count() (int64, error) {
	var count int64
	var entity T
	err := b.db.Model(&entity).Count(&count).Error
	return count, err
}

func (b *BaseDB[T]) IterateFindByColumn(column string, desc bool, limit int, handler func(entities []*T) error, wait time.Duration) error {
	var offset = 0
	for {
		entities, err := b.FindByColumn(column, desc, limit, offset)
		if err != nil {
			return err
		}
		if len(entities) == 0 {
			break
		}
		if err = handler(entities); err != nil {
			return err
		}
		offset += limit
		if wait > 0 {
			time.Sleep(wait)
		}
	}
	return nil
}

// DropTable drops the table, used for testing
func (b *BaseDB[T]) DropTable() error {
	var entity T
	return b.db.Migrator().DropTable(&entity)
}
