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
	err := db.Order("effective_on desc").
		Where("os = ?", os).
		Where("effective_on <= ?", time.Now()).
		First(&appVersion).Error

	return &appVersion, err
}

func (rp *utilsRepo) SaveUpload(modelInstance *models.Upload) error {
	return saveDbModel(modelInstance)
}

func (rp *utilsRepo) SaveConfig(modelInstance *models.Config) error {
	return saveDbModel(modelInstance)
}

func (rp *utilsRepo) SaveAppVersion(modelInstance *models.AppVersion) error {
	return saveDbModel(modelInstance)
}

func (rp *utilsRepo) SaveTraffic(modelInstance *models.Traffic) error {
	return saveDbModel(modelInstance)
}
