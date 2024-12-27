package repositories

import (
	"strconv"

	"github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/utils"

	"gorm.io/gorm"
)

type folderRepo struct {
	mainDBCon *gorm.DB
	readDBCon *gorm.DB
	rd        *redisRepo
}

// Get user folders
func (rp *folderRepo) GetUserFolders(
	userId uint64,
) (
	*[]models.Folder, error,
) {
	if userId <= 0 {
		return nil, errors.ApiErrorParamInvalid
	}

	db := rp.readDBCon
	folders := []models.Folder{}
	if err := db.Model(&models.Folder{}).Preload("Invitees").Find(&folders, "owner_id = ?", userId).
		Error; err != nil {
		return nil, err
	}
	return &folders, nil
}

// Get folder by field name
func (rp *folderRepo) GetFolderByFieldName(
	fieldName, fieldValue string,
) (
	*models.Folder, error,
) {
	db := rp.readDBCon
	var folder models.Folder
	if err := db.Where(fieldName+" = ?", fieldValue).
		First(&folder).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

// Save folder
func (rp *folderRepo) SaveFolder(
	user *models.User, folderId uint64, updateData *map[string]any, checkUserAuthentication bool,
) (
	*models.Folder, error,
) {
	db := rp.mainDBCon
	if folderId <= 0 {
		folder := &models.Folder{}
		if err := models.PopulateModelFromMap(folder, *updateData); err != nil {
			return nil, err
		}

		folder.OwnerID = user.ID
		if err := db.Create(folder).Error; err != nil {
			return nil, err
		}
		return folder, nil
	} else {
		folder, err := rp.GetFolderByFieldName("id", strconv.Itoa(int(folderId)))
		if err != nil {
			return nil, err
		}
		if checkUserAuthentication && !rp.HasFolderUpdatePermission(user, folder) {
			return nil, errors.ApiErrorPermissionDenied
		}
		if err := models.PopulateModelFromMap(folder, *updateData); err != nil {
			return nil, err
		}
		if err := SaveDbModel(folder); err != nil {
			return nil, err
		}
		return folder, nil
	}
}

// Sync inviteeIDs
func (rp *folderRepo) SyncInviteeIDs(
	user *models.User, folder *models.Folder, newInviteeIds *[]uint64,
) error {
	if !rp.HasFolderUpdatePermission(user, folder) {
		return errors.ApiErrorPermissionDenied
	}

	db := rp.mainDBCon
	currentInviteeIds := []uint64{}
	db.Model(&models.FolderInvitee{}).
		Where("folder_id = ?", folder.ID).
		Pluck("invitee_id", &currentInviteeIds)

	needsDeleteInviteeIds := utils.SliceLeftExcludeRight(&currentInviteeIds, newInviteeIds)
	needsInsertInviteeIds := utils.SliceLeftExcludeRight(newInviteeIds, &currentInviteeIds)

	if len(*needsDeleteInviteeIds) > 0 {
		if err := db.
			Where("folder_id = ?", folder.ID).
			Where("invitee_id in ?", *needsDeleteInviteeIds).
			Delete(&models.FolderInvitee{}).Error; err != nil {
			return err
		}
	}
	if len(*needsInsertInviteeIds) > 0 {
		invitees := make([]models.FolderInvitee, 0)
		for _, item := range *needsInsertInviteeIds {
			invitees = append(invitees,
				models.FolderInvitee{FolderID: folder.ID, InviteeID: item})
		}
		if err := db.Create(invitees).Error; err != nil {
			return err
		}
	}

	return nil
}

// Delete target folder with folderId
func (rp *folderRepo) DeleteFolder(
	user *models.User, folderId uint64,
) error {
	if folderId <= 0 {
		return errors.ApiErrorParamInvalid
	}

	db := rp.mainDBCon
	var folder models.Folder
	if err := QueryDbModelByPrimaryId(&folder, folderId); err != nil {
		return err
	}

	if !rp.HasFolderUpdatePermission(user, &folder) {
		return errors.ApiErrorPermissionDenied
	}

	return db.Delete(folder).Error
}

// Checks if user have permission to update the folder
// User can update if or(user.IsAdmin, user.IsOwner)
func (rp *folderRepo) HasFolderUpdatePermission(
	user *models.User, folder *models.Folder,
) bool {
	if user == nil || folder == nil {
		return false
	}
	if user.IsAdmin() {
		return true
	}
	if folder.OwnerID == user.ID {
		return true
	}

	db := rp.readDBCon
	count := int64(0)
	db.Model(&models.FolderInvitee{}).
		Where("folder_id", folder.ID).
		Where("invitee_id", user.ID).
		Count(&count)
	return false
}

// Checks if user have permission to post to the folder
// User can post if or(user.IsAdmin, user.IsOwner, user.IsFolderInvitee)
func (rp *folderRepo) HasFolderPostPermission(
	user *models.User, folder *models.Folder,
) bool {
	if rp.HasFolderUpdatePermission(user, folder) {
		return true
	}

	db := rp.readDBCon
	count := int64(0)
	db.Model(&models.FolderInvitee{}).
		Where("folder_id", folder.ID).
		Where("invitee_id", user.ID).
		Count(&count)

	return count > 0
}

func (rp *folderRepo) FindInvitees(folder *models.Folder) []*models.UserMember {
	db := rp.readDBCon
	var invitees []*models.UserMember
	db.Model(&folder).Association("Invitees").Find(&invitees)
	return invitees
}
