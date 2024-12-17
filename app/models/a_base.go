package models

import (
	"reflect"
	"time"

	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BaseDbModel struct {
	ID        int64     `gorm:"column:id; primary_key; not null" json:"id"`
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
