package models

import "time"

type FolderFileUpload struct {
	BaseDbModel
	UserId   int64     `gorm:"column:user_id" json:"userId"`
	Disk     string    `gorm:"column:disk" json:"-"`
	FileName string    `gorm:"column:file_name" json:"fileName"`
	FileType string    `gorm:"column:file_type; type:enum('image', 'video', 'pdf', 'doc', 'excel', 'others')" json:"fileType"`
	FilePath string    `gorm:"column:file_path" json:"filePath"`
	ShotAt   time.Time `gorm:"column:shot_at" json:"shotAt"`
}

func (FolderFileUpload) TableName() string {
	return "folder_file_uploads"
}