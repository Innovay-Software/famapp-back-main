package repositories

import (
	"github.com/innovay-software/famapp-main/app/models"
	"gorm.io/gorm"
)

type cloudUploadRepo struct {
	mainDBCon *gorm.DB
	readDBCon *gorm.DB
}

// const testingId = 7700

// Get a list of folder files to upload to HW OBS
// Condition: hw_original_file_path is ""
func (cu *cloudUploadRepo) GetHwObsCandidates(
	limit int,
) (
	*[]models.FolderFile, error,
) {
	// Use mainDB to get latest data
	db := cu.mainDBCon
	folderFiles := []models.FolderFile{}
	result := db.Model(&models.FolderFile{}).
		Scopes(folderFileScopeActive, folderFileScopeNotDownloading).
		// Where("id = ?", testingId).
		Where("hw_original_file_path = ?", "").
		Order("id asc").
		Limit(limit).
		Find(&folderFiles)

	return &folderFiles, result.Error
}

// Get a list of folder files to upload to Google Storage
// Conditions:
//
//	Uploaded to HW OBS AND
//	Uploaded to Google Drive AND
//	any one of (google_original_file_path, google_file_path, google_thumbnail_path) is empty
//
// That is, if any one the three google storage paths are missing, re-upload everything to google storage
func (cu *cloudUploadRepo) GetGoogleStorageCandidates(
	limit int,
) (
	*[]models.FolderFile, error,
) {
	// Use mainDB to get latest data
	db := cu.mainDBCon
	folderFiles := []models.FolderFile{}
	result := db.Model(&models.FolderFile{}).
		Scopes(folderFileScopeActive, folderFileScopeNotDownloading).
		// Where("id = ?", testingId).
		Where("hw_original_file_path not in ?", []string{"processing", ""}).
		Where("google_drive_file_id not in ?", []string{"processing", ""}).
		Where(
			db.Where("google_original_file_path = ?", "").
				Or("google_file_path = ?", "").
				Or("google_thumbnail_path = ?", ""),
		).
		Order("id asc").
		Limit(limit).
		Find(&folderFiles)
	return &folderFiles, result.Error
}

// Get a list of folder files to upload to Google Storage
func (cu *cloudUploadRepo) GetGoogleDriveCandidates(
	limit int,
) (
	*[]models.FolderFile, error,
) {
	// Use mainDB to get latest data
	db := cu.mainDBCon
	folderFiles := []models.FolderFile{}
	result := db.Model(&models.FolderFile{}).
		Scopes(folderFileScopeActive, folderFileScopeNotDownloading).
		// Where("id = ?", testingId).
		Where("google_drive_file_id = ?", "").
		Order("id asc").
		Limit(limit).
		Find(&folderFiles)
	return &folderFiles, result.Error
}

func (cu *cloudUploadRepo) BatchUploadFolderFile(
	fieldName, fieldValue string, ids []int,
) error {
	// Use mainDB to get latest data
	db := cu.mainDBCon
	result := db.Model(&models.FolderFile{}).
		Where("id in ?", ids).
		Update(fieldName, fieldValue)
	return result.Error
}
