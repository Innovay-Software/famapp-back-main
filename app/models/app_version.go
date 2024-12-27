package models

import "time"

type AppVersion struct {
	BaseDbModel
	OS          string    `gorm:"column:os; type:enum('ios', 'android')" json:"os"`
	Version     string    `gorm:"column:version" json:"version"`
	Title       string    `gorm:"column:title" json:"title"`
	Content     JSONB     `gorm:"column:content; type:jsonb; default:'{}'" json:"content"`
	Link        string    `gorm:"column:link" json:"link"`
	DownloadUrl string    `gorm:"column:download_url" json:"download_url"`
	EffectiveOn time.Time `gorm:"column:effective_on" json:"effective_on"`
}

func (AppVersion) TableName() string {
	return "app_versions"
}
