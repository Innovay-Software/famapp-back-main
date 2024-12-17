package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/innovay-software/famapp-main/app/utils"
)

type Upload struct {
	BaseDbModel
	UserID   int64     `gorm:"column:user_id" json:"userId"`
	Disk     string    `gorm:"column:disk" json:"disk"`
	FileName string    `gorm:"column:file_name" json:"fileName"`
	FileType string    `gorm:"column:file_type" json:"fileType"`
	FilePath string    `gorm:"column:file_path" json:"filePath"`
	ShotAt   time.Time `gorm:"column:shot_at" json:"shotAt"`
	FileUrl  string    `gorm:"-" json:"fileUrl"`
}

func (Upload) TableName() string {
	return "uploads"
}

func (model *Upload) PopulateDerivedFields() {
	model.FileUrl = utils.GetUrlPath(model.Disk, model.FilePath)
	fmt.Println("Calculate Drived Field: ", model.FileUrl)
}

func (model Upload) MarshalJSON() ([]byte, error) {
	model.PopulateDerivedFields()
	type Upload2 Upload
	return json.Marshal(Upload2(model))
}
