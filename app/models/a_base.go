package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BaseDbModel struct {
	ID        uint64    `gorm:"column:id; primary_key; not null" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// func (*BaseDbModel) GetFileUrl(disk, filepath string) string {
// 	if filepath != "" && !strings.HasPrefix(filepath, "http") {
// 		filepath = utils.GetUrlPath(disk, filepath)
// 	}
// 	return filepath
// }

type BaseModelSoftDelete struct {
	BaseDbModel
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`
}

func PopulateModelFromMap[T any](modelInstance *T, data map[string]any) error {
	err := mapstructure.Decode(data, modelInstance)
	if err != nil {
		utils.Log(reflect.TypeOf(data["metadata"]))
		utils.Log("Populate model from map filed:", err)
		return err
	}
	return nil
}

// JSONB Interface for JSONB Field of yourTableName Table
type JSONB map[string]interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
