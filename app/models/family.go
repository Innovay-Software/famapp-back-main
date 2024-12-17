package models

type Family struct {
	BaseDbModel
	Title  string `gorm:"column:title" json:"title"`
	Status bool  `gorm:"column:status" json:"status"`
}

func (Family) TableName() string {
	return "families"
}
