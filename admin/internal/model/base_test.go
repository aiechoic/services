package model_test

import (
	"fmt"
	"github.com/aiechoic/services/admin/internal/model"
	"github.com/aiechoic/services/database/gorm"
	"github.com/aiechoic/services/ioc"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestEntity struct {
	ID    int
	Name  string
	Value string
}

func setupBaseDB(t *testing.T) (*model.BaseDB[TestEntity], func()) {
	c := ioc.NewContainer()
	err := c.LoadConfig(configPath, ioc.ConfigEnvTest)
	if err != nil {
		t.Fatal(err)
	}
	db := gorm.GetGormDB(c)
	baseDB, err := model.NewBaseDB[TestEntity](db)
	if err != nil {
		t.Fatal(err)
	}
	deferFunc := func() {
		err = baseDB.DropTable()
		assert.NoError(t, err)
	}
	return baseDB, deferFunc
}

func TestBaseDB_Create(t *testing.T) {
	baseDB, closer := setupBaseDB(t)
	defer closer()

	entity := &TestEntity{Name: "name", Value: "value"}
	err := baseDB.Create(entity)
	assert.NoError(t, err)
	need := &TestEntity{ID: 1, Name: "name", Value: "value"}
	assert.Equal(t, need, entity)
}

func TestBaseDB_GetByColumn(t *testing.T) {
	baseDB, closer := setupBaseDB(t)
	defer closer()

	entity := &TestEntity{Name: "name", Value: "value"}
	_ = baseDB.Create(entity)
	queryEntity, err := baseDB.GetByColumn("id", entity.ID)
	assert.NoError(t, err)
	need := &TestEntity{ID: 1, Name: "name", Value: "value"}
	assert.Equal(t, need, queryEntity)
}

func TestBaseDB_FullUpdateByColumn(t *testing.T) {
	baseDB, closer := setupBaseDB(t)
	defer closer()

	entity := &TestEntity{Name: "name", Value: "value"}
	_ = baseDB.Create(entity)
	entity.Name = "" // save the empty value
	saveEntity, err := baseDB.FullUpdateByColumn(entity, "id", entity.ID)
	assert.NoError(t, err)
	need := &TestEntity{ID: 1, Name: "", Value: "value"}
	assert.Equal(t, need, saveEntity)
}

func TestBaseDB_PartialUpdateByColumn(t *testing.T) {
	baseDB, closer := setupBaseDB(t)
	defer closer()

	entity := &TestEntity{Name: "name", Value: "value"}
	_ = baseDB.Create(entity)
	entity.Name = "new_name"
	entity.Value = "" // empty value will not be updated
	updateEntity, err := baseDB.PartialUpdateByColumn(entity, "id", entity.ID)
	assert.NoError(t, err)
	need := &TestEntity{ID: 1, Name: "new_name", Value: "value"}
	assert.Equal(t, need, updateEntity)
}

func TestBaseDB_DeleteByColumn(t *testing.T) {
	baseDB, closer := setupBaseDB(t)
	defer closer()

	entity := &TestEntity{Name: "name", Value: "value"}
	_ = baseDB.Create(entity)
	gotRowsAffected := baseDB.DeleteByColumn("id", entity.ID)
	assert.Equal(t, int64(1), gotRowsAffected)
}

func TestBaseDB_FindByColumn(t *testing.T) {
	baseDB, closer := setupBaseDB(t)
	defer closer()

	entity := &TestEntity{Name: "name", Value: "value"}
	_ = baseDB.Create(entity)
	entities, err := baseDB.FindByColumn("id", false, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, entities, 1)
}

func TestBaseDB_Count(t *testing.T) {
	baseDB, closer := setupBaseDB(t)
	defer closer()

	entity := &TestEntity{Name: "name", Value: "value"}
	_ = baseDB.Create(entity)
	count, err := baseDB.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestBaseDB_IterateFindByColumn(t *testing.T) {
	baseDB, closer := setupBaseDB(t)
	defer closer()

	// Create multiple entities
	for i := 1; i <= 5; i++ {
		entity := &TestEntity{Name: fmt.Sprintf("name%d", i), Value: fmt.Sprintf("value%d", i)}
		_ = baseDB.Create(entity)
	}

	// Handler to collect entities
	var collectedEntities []*TestEntity
	handler := func(entities []*TestEntity) error {
		collectedEntities = append(collectedEntities, entities...)
		return nil
	}

	// Call IterateFindByColumn
	err := baseDB.IterateFindByColumn("id", false, 2, handler, 0)
	assert.NoError(t, err)
	assert.Len(t, collectedEntities, 5)

	// Verify the collected entities
	for i, entity := range collectedEntities {
		assert.Equal(t, fmt.Sprintf("name%d", i+1), entity.Name)
		assert.Equal(t, fmt.Sprintf("value%d", i+1), entity.Value)
	}
}
