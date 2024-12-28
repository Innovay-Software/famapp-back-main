package repositories

import (
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/utils"
	"gorm.io/gorm"
)

type jobRepo struct {
	mainDBCon *gorm.DB
}

func (rp *jobRepo) GetUploadToGoogleDriveCandidate() (*models.FolderFile, error) {
	db := rp.mainDBCon

	var folderFile models.FolderFile
	err1 := db.Scopes(folderFileScopeActive, folderFileScopeReadyForGoogleUpload).
		Where("google_drive_file_id IS NULL").
		Or("google_drive_file_id = ?", "").
		First(&folderFile).Error

	if err1 != nil {
		utils.LogError("Error getting a folderFileCandidate:", err1)
		return nil, err1
	}

	err2 := db.Scopes(folderFileScopeActive, folderFileScopeReadyForGoogleUpload).
		Where("id = ?", folderFile.ID).
		Where(
			db.Where("google_drive_file_id IS NULL").
				Or("google_drive_file_id = ?", "")).
		Update("google_drive_file_id", "processing").Error

	if err2 != nil {
		utils.LogError("Error updating folderFile to processing:", err2)
		return nil, err2
	}

	return &folderFile, nil
}

func (rp *jobRepo) GetFolderFile(folderFileId uint64) (*models.FolderFile, error) {
	db := rp.mainDBCon
	var folderFile models.FolderFile
	if err := db.Where("id = ?", folderFileId).First(&folderFile).Error; err != nil {
		return nil, err
	}
	return &folderFile, nil
}