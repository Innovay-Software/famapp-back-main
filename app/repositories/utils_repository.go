package repositories

import (
	"time"

	"github.com/innovay-software/famapp-main/app/models"

	"gorm.io/gorm"
)

type utilsRepo struct {
	mainDBCon *gorm.DB
	readDBCon *gorm.DB
}

// Get user folders
func (rp *utilsRepo) GetLatestAppVersion(
	os string,
) (
	*models.AppVersion, error,
) {
	db := rp.readDBCon
	var appVersion models.AppVersion
	err := db.Order("effective_date desc").
		Where("os = ?", os).
		Where("effective_date <= ?", time.Now()).
		First(&appVersion).Error

	return &appVersion, err
}
