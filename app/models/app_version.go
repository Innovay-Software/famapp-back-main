package models

import "time"

type AppVersion struct {
	BaseDbModel
	OS            string    `gorm:"column:os; type:enum('ios', 'android')" json:"os"`
	Version       string    `gorm:"column:version" json:"version"`
	Title         string    `gorm:"column:title" json:"title"`
	Content       string    `gorm:"column:content" json:"content"`
	Link          string    `gorm:"column:link" json:"link"`
	File          string    `gorm:"column:file" json:"file"`
	EffectiveDate time.Time `gorm:"column:effective_date" json:"effective_date"`
}

func (AppVersion) TableName() string {
	return "app_versions"
}
