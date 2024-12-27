package models

import (
	"encoding/json"
	"time"

	"github.com/innovay-software/famapp-main/app/utils"
)

type Upload struct {
	BaseDbModel
	UserID   uint64    `gorm:"column:user_id" json:"userId"`
	Disk     string    `gorm:"column:disk" json:"disk"`
	FileName string    `gorm:"column:file_name" json:"fileName"`
	FileType string    `gorm:"column:file_type" json:"fileType"`
	FilePath string    `gorm:"column:file_path" json:"filePath"`
	TakenOn  time.Time `gorm:"column:taken_on" json:"takenOn"`
	FileUrl  string    `gorm:"-" json:"fileUrl"`
}

func (Upload) TableName() string {
	return "uploads"
}

func (model *Upload) PopulateDerivedFields() {
	model.FileUrl = utils.GetUrlPath(model.Disk, model.FilePath)
}

func (model Upload) MarshalJSON() ([]byte, error) {
	model.PopulateDerivedFields()
	type Upload2 Upload
	return json.Marshal(Upload2(model))
}
