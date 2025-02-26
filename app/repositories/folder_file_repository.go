package repositories

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"time"

	"github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/utils"
)

// Folder file repository uses the same folderRepo with folder repository

const (
	pageSizeForOlderRecords = 50
	pageSizeForNewerRecords = 30
)

func (rp *folderRepo) CountFilesForTargetUnixTimestamp(timestamp int64) uint64 {
	db := rp.readDBCon
	fromDateTime := time.Unix(timestamp, 0)
	toDateTime := time.Unix(timestamp+1, 0)
	fileCount := int64(0)
	db.Model(&models.FolderFile{}).
		Where("taken_on >= ?", fromDateTime).
		Where("taken_on < ?", toDateTime).
		Count(&fileCount)

	if fileCount < 999999 {
		return uint64(fileCount)
	}

	// if there are more than 1 mil files at target timestamp, increase timestamp by 1 second
	return rp.CountFilesForTargetUnixTimestamp(timestamp + 1)
}

// Save folder file instance
func (rp *folderRepo) CreateFolderFile(
	folderFile *models.FolderFile,
) error {
	err := saveDbModel(folderFile)
	if err != nil {
		return err
	}
	MessageQueueRepoIns.sendFolderFileToFolderFileProcessingQueue(folderFile)
	return nil
}

// Delete folder file record
func (rp *folderRepo) DeleteFolderFile(
	user *models.User, folderFileId uint64,
) error {
	folderFileIdString := strconv.Itoa(int(folderFileId))
	folderFile, err := rp.GetFolderFileByFieldName(user, "id", folderFileIdString)
	if err != nil {
		return err
	}
	folderFile.FolderID = 0
	return saveDbModel(folderFile)
}

func (rp *folderRepo) SaveFolderFileModel(folderFile *models.FolderFile) error {
	if takenOnDateTime, exists := folderFile.Metadata["taken_on_date_time"]; exists {
		utils.Log("TakenOnDateTime type=", reflect.TypeOf(takenOnDateTime))
		// if takenOnDateTime, ok := takenOnDateTime.(time.Time); ok {

		// }
	}
	return saveDbModel(folderFile)
}

// Update folder file
func (rp *folderRepo) UpdateFolderFile(
	user *models.User, folderFileId, newFolderId uint64, remark string, isPrivate bool,
) (
	*models.FolderFile, error,
) {
	if folderFileId <= 0 {
		return nil, errors.ApiErrorParamInvalid
	}

	folderFile, err := rp.GetFolderFileByFieldName(user, "id", strconv.Itoa(int(folderFileId)))
	if err != nil {
		return nil, err
	}
	if folderFile.OwnerID != user.ID && !user.IsAdmin() {
		folder, err := rp.GetFolderByFieldName("id", strconv.Itoa(int(folderFile.FolderID)))
		if err != nil {
			return nil, err
		}
		if folder.OwnerID != user.ID {
			return nil, errors.ApiErrorPermissionDenied
		}
	}
	folderFile.Remark = remark
	folderFile.IsPrivate = isPrivate
	if newFolderId >= 0 {
		folderFile.FolderID = uint64(newFolderId)
	}
	saveToDbErr := saveDbModel(folderFile)
	return folderFile, saveToDbErr
}

// Query folder file
func (rp *folderRepo) GetFolderFileByFieldName(
	user *models.User, fieldName string, fieldValue string,
) (
	*models.FolderFile, error,
) {
	db := rp.readDBCon
	var folderFile models.FolderFile
	if err := db.Where(fieldName+" = ?", fieldValue).First(&folderFile).Error; err != nil {
		return nil, err
	}
	if folderFile.OwnerID != user.ID {
		// not owner, check if user has update permission
		var folder models.Folder
		if err := db.First(&folder, folderFile.FolderID).Error; err != nil {
			return nil, err
		}

		if !rp.HasFolderUpdatePermission(user, &folder) {
			return nil, errors.ApiErrorPermissionDenied
		}
	}
	return &folderFile, nil
}

// Query for folder files with target MD5 value
func (rp *folderRepo) GetActiveFolderFileWithMd5(
	folderId uint64, md5Value string,
) (
	*models.FolderFile, error,
) {
	if md5Value == "" {
		return nil, errors.ApiError{Message: "Missing md5 value"}
	}

	db := rp.readDBCon
	md5Value += "%"
	var folderFile models.FolderFile
	err := db.Scopes(folderFileScopeActive).
		Where("folder_id = ?", folderId).
		Where("file_name like ?", md5Value).
		First(&folderFile).Error

	if err != nil {
		return nil, err
	}
	return &folderFile, nil
}

// Get latest
func (rp *folderRepo) GetLatestFolderFileForDateTimeSecond(
	folderId uint64, targetDatetime time.Time,
) (
	*models.FolderFile, error,
) {
	if folderId <= 0 {
		return nil, errors.ApiErrorParamInvalid
	}
	db := rp.readDBCon

	var folderFile models.FolderFile
	timeString := targetDatetime.Format(time.DateTime)

	if err := db.Where("folder_id", folderId).
		Where("taken_on like ", timeString+"%").
		Order("taken_on desc").
		First(&folderFile).Error; err != nil {
		return nil, err
	}

	return &folderFile, nil
}

// Get past records before a certain point
func (rp *folderRepo) GetFolderFilesBeforeTakenOn(
	user *models.User, folderId uint64, pivotDate string, beforeMicroTimestamp int64,
) (
	*[]models.FolderFile, *models.Folder, bool, error,
) {
	if folderId <= 0 {
		return nil, nil, false, errors.ApiErrorParamInvalid
	}
	folder, err := rp.GetFolderByFieldName("id", strconv.Itoa(int(folderId)))
	if err != nil {
		return nil, nil, false, err
	}

	folderFiles, hasMore, err := rp.getFolderFiles(
		user, folder, pageSizeForOlderRecords, "<", "desc", pivotDate, beforeMicroTimestamp,
	)
	return folderFiles, folder, hasMore, err
}

// Get newer records after a certain point
func (rp *folderRepo) GetFolderFilesAfterTakenOn(
	user *models.User, folderId uint64, pivotDate string, afterMicroTimestamp int64,
) (
	*[]models.FolderFile, *models.Folder, bool, error,
) {
	if folderId <= 0 {
		return nil, nil, false, errors.ApiErrorParamInvalid
	}

	folder, err := rp.GetFolderByFieldName("id", strconv.Itoa(int(folderId)))
	if err != nil {
		return nil, nil, false, err
	}

	folderFiles, hasMore, err := rp.getFolderFiles(
		user, folder, pageSizeForNewerRecords, ">", "asc", pivotDate, afterMicroTimestamp,
	)
	return folderFiles, folder, hasMore, err
}

// Returns a slice of files and has more bool
func (rp *folderRepo) getFolderFiles(
	user *models.User, folder *models.Folder, pageSize int,
	sign, order, pivotDate string, pivotMicroTimestampTakenOn int64,
) (
	*[]models.FolderFile, bool, error,
) {
	db := rp.readDBCon

	if pivotDate == "-" {
		pivotDate = ""
	}
	if pivotDate != "" {
		if _, err := time.Parse(time.DateOnly, pivotDate); err != nil {
			pivotDate = ""
		}
	}

	var pivotTakenOnTime *time.Time
	if pivotMicroTimestampTakenOn > 0 {
		seconds := pivotMicroTimestampTakenOn / 1000000
		micro := pivotMicroTimestampTakenOn % 1000000
		t := time.Unix(seconds, micro*1000)
		pivotTakenOnTime = &t
	}

	var folderFiles []models.FolderFile

	if folder.OwnerID != user.ID && !user.IsAdmin() {
		// user is not owner nor admin, check if is an invitee of this folder
		inviteeRecord := models.FolderInvitee{}
		if err := db.Where("folder_id = ?", folder.ID).
			Where("invitee_id = ?", user.ID).
			First(&inviteeRecord).Error; err != nil {
			return &folderFiles, false, err
		}
	}

	query := db.Limit(pageSize).
		Model(&models.FolderFile{}).
		Where("folder_id = ?", folder.ID).
		Where("file_type in ?", []string{"image", "video"})

	if pivotMicroTimestampTakenOn > 0 {
		query = query.Where("taken_on "+sign+" ?", *pivotTakenOnTime)
	}
	if pivotDate != "" {
		date, err := time.Parse("2006-01-02", pivotDate)
		if err == nil {
			date = date.AddDate(0, 0, 1)
			query = query.Where("taken_on "+sign+"= ?", date)
		}
	}
	query.Order("taken_on " + order).Find(&folderFiles)

	if sign == ">=" && len(folderFiles) > 0 {
		utils.ReverseSliceInPlace(&folderFiles)
	}

	hasMore := len(folderFiles) >= pageSize
	filteredFolderFiles := filterOutFolderFilesWithoutViewPermission(user, folder, &folderFiles)
	folderFiles = *filteredFolderFiles

	return &folderFiles, hasMore, nil
}

func filterOutFolderFilesWithoutViewPermission(
	user *models.User, folder *models.Folder, folderFiles *[]models.FolderFile,
) *[]models.FolderFile {
	if folder.OwnerID == user.ID || user.IsAdmin() {
		return folderFiles
	}

	filteredFiles := []models.FolderFile{}
	for _, item := range *folderFiles {
		if item.OwnerID == user.ID || !item.IsPrivate {
			filteredFiles = append(filteredFiles, item)
		}
	}
	return &filteredFiles
}

func (rp *folderRepo) MoveFolderFiles(
	user *models.User, folderFileIds *[]uint64, newFolderId uint64, processLimit int,
) error {
	if folderFileIds == nil || len(*folderFileIds) == 0 {
		return fmt.Errorf("missing folder file ids")
	}
	db := mainDBCon
	var folder models.Folder
	if err := db.First(&folder, newFolderId).Error; err != nil {
		return err
	}

	*folderFileIds = (*folderFileIds)[:processLimit]

	if !user.IsAdmin() {
		var folderFiles []models.FolderFile
		if err := db.Limit(processLimit).Model(
			&models.FolderFile{},
		).Where("id in ?", *folderFileIds).Find(&folderFiles).Error; err != nil {
			return err
		}

		filteredFolderIds := []uint64{}
		folderIdMap := map[uint64][]uint64{}
		for _, ff := range folderFiles {
			if ff.OwnerID == user.ID {
				filteredFolderIds = append(filteredFolderIds, ff.ID)
			} else {
				if folderIdMap[ff.FolderID] == nil {
					folderIdMap[ff.FolderID] = []uint64{}
				}
				folderIdMap[ff.FolderID] = append(
					folderIdMap[ff.FolderID],
					ff.ID,
				)
			}
		}

		folderIds := []uint64{}
		for k := range folderIdMap {
			folderIds = append(folderIds, k)
		}
		var folders []models.Folder
		if err := db.Model(&models.Folder{}).Where(
			"id in ?", folderIds,
		).Find(&folders).Error; err == nil {
			for _, f := range folders {
				if f.OwnerID == user.ID {
					filteredFolderIds = append(filteredFolderIds, folderIdMap[f.ID]...)
				}
			}
		}
		folderFileIds = &filteredFolderIds
	}

	return db.Model(&models.FolderFile{}).Where(
		"id in ?", *folderFileIds,
	).Update("folder_id", newFolderId).Error
}

func (rp *folderRepo) RescheduleFolderFiles(
	user *models.User, folderFileIds *[]uint64, newTimestampInSeconds int64, processLimit int,
) error {
	if folderFileIds == nil || len(*folderFileIds) == 0 {
		return fmt.Errorf("missing folder file ids")
	}

	db := rp.mainDBCon
	*folderFileIds = (*folderFileIds)[:processLimit]
	var folderFiles []models.FolderFile
	if err := db.Limit(processLimit).
		Model(&models.FolderFile{}).
		Where("id in ?", *folderFileIds).Find(&folderFiles).Error; err != nil {
		return err
	}

	deniedFolderFileIds := []uint64{}
	if !user.IsAdmin() {
		// Check the owners of the candidate folder files
		folderIdMap := map[uint64][]uint64{}
		for _, ff := range folderFiles {
			if ff.OwnerID != user.ID {
				if folderIdMap[ff.FolderID] == nil {
					folderIdMap[ff.FolderID] = []uint64{}
				}
				folderIdMap[ff.FolderID] = append(
					folderIdMap[ff.FolderID],
					ff.ID,
				)
			}
		}

		folderIds := []uint64{}
		for k := range folderIdMap {
			folderIds = append(folderIds, k)
		}
		var folders []models.Folder
		if err := db.Model(&models.Folder{}).
			Where("id in ?", folderIds).
			Find(&folders).Error; err == nil {

			for _, f := range folders {
				if f.OwnerID != user.ID {
					deniedFolderFileIds = append(
						deniedFolderFileIds,
						folderIdMap[f.ID]...,
					)
				}
			}
		}
	}

	// find the latest folder file in that second
	secondWindowStart := time.Unix(newTimestampInSeconds, 0)
	secondWindowEnd := time.Unix(newTimestampInSeconds+1, 0)
	targetMilliSecond := 1
	var latestFolderFile models.FolderFile
	if err := db.Model(&models.FolderFile{}).
		Where("taken_on >= ?", secondWindowStart).
		Where("taken_on < ?", secondWindowEnd).
		Order("taken_on desc").
		First(&latestFolderFile).Error; err != nil {

		targetMilliSecond = latestFolderFile.TakenOn.Nanosecond() / 1000
	}

	for _, ff := range folderFiles {
		if slices.Contains(deniedFolderFileIds, ff.ID) {
			continue
		}
		ff.TakenOn = time.Unix(newTimestampInSeconds, int64(targetMilliSecond)*1000)
		saveDbModel(&ff)
		targetMilliSecond++
	}

	return nil
}
