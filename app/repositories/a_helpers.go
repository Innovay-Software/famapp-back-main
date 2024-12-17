package repositories

import (
	"fmt"
	"reflect"

	"gorm.io/gorm/schema"
)

// A wrapper around db.Save method. This wrapper is created to solve the common error where
// a model instance is passed to Save function instead of a pointer, resulting in a panic
// This wrapper will check if the input param is a pointer to a DbModel before calling
// the Save function to try to save it


// Save will create the record if the primary key is not provided
func CreateDbModel(model schema.Tabler) error {
	modelType := reflect.TypeOf(model).String()
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("a pointer to a model instance is required, not %s", modelType)
	}
	return mainDBCon.Create(model).Error
}

// Save will create the record if the primary key is not provided
func SaveDbModel(model schema.Tabler) error {
	modelType := reflect.TypeOf(model).String()
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("a pointer to a model instance is required, not %s", modelType)
	}
	return mainDBCon.Save(model).Error
}

// Delele model
func DeleteDbModel(model schema.Tabler) error {
	modelType := reflect.TypeOf(model).String()
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("a pointer to a model instance is required, not %s", modelType)
	}
	return mainDBCon.Delete(model).Error
}

// Query by primary key
func QueryDbModelByPrimaryId(model schema.Tabler, id int64) error {
	modelType := reflect.TypeOf(model).String()
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("a pointer to a model instance is required, not %s", modelType)
	}
	if id <= 0 {
		return fmt.Errorf("record not found")
	}
	return readDBCon.First(model, id).Error
}
