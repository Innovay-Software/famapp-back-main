package repositories

import (
	"time"

	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/utils"
	"gorm.io/gorm"
)

type configRepo struct {
	mainDBCon *gorm.DB
	readDBCon *gorm.DB
	rd        *redisRepo
}

const (
	uploadTimeKey = "UploadTime"
)

// Read config from DB
func (cr *configRepo) GetConfig(configKey string) *models.Config {
	db := cr.readDBCon
	var config models.Config
	if db.First(&config, "config_key = ?", configKey).Error != nil {
		return nil
	}
	return &config
}

// Count number of config records in database
func (cr *configRepo) CountConfigRecords() int64 {
	db := cr.readDBCon
	count := int64(0)
	db.Model(&models.Config{}).Count(&count)
	return count
}

// Set latest file upload time
func (cr *configRepo) UpdateUploadTime() error {
	configVal := time.Now().UTC().Format(time.RFC3339Nano)
	db := cr.mainDBCon

	var config models.Config
	if db.First(&config, "config_key = ?", uploadTimeKey).Error != nil {
		// Record not found, create a new one
		config = models.Config{
			Title:       uploadTimeKey,
			ConfigKey:   uploadTimeKey,
			ConfigValue: configVal,
		}
		result := db.Create(&config)
		return result.Error
	} else {
		config.ConfigValue = configVal
		return saveDbModel(&config)
	}
}

func (cr *configRepo) GetUploadTime() *time.Time {
	db := cr.readDBCon
	now := time.Now().UTC()
	var config models.Config
	if err := db.Where("config_key = ?", uploadTimeKey).
		First(&config).Error; err != nil {
		utils.Log("Can't get config:", err)
		return &now
	}
	uploadTime, err := time.Parse(time.RFC3339Nano, config.ConfigValue)
	if err != nil {
		utils.Log("Can't parse config:", err)
		return &now
	}
	return &uploadTime
}
