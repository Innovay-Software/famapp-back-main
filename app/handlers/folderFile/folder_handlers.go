package folderFile

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
	"github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
)

func SaveFolderHandler(
	c *gin.Context, user *models.User, folderId, ownerId, parentId int64,
	title, cover, folderType string, isDefault, isPrivate bool, metadata *map[string]any,
	inviteeIds *[]int64,
) (
	dto.ApiResponse, error,
) {
	var folder models.Folder
	if err := repositories.QueryDbModelByPrimaryId(
		&folder, folderId,
	); err != nil {
		folder = models.Folder{OwnerID: user.ID}
	}

	if !repositories.FolderRepoIns.HasFolderUpdatePermission(user, &folder) {
		return nil, errors.ApiErrorPermissionDenied
	}

	folder.OwnerID = ownerId
	folder.ParentID = parentId
	folder.Title = title
	folder.Cover = cover
	folder.Type = folderType
	folder.IsDefault = isDefault
	folder.IsPrivate = isPrivate
	if metadata != nil {
		folder.Metadata = *metadata
	}

	if err := repositories.SaveDbModel(&folder); err != nil {
		return nil, err
	}

	// Update invitees
	if !slices.Contains(*inviteeIds, ownerId) {
		*inviteeIds = append(*inviteeIds, ownerId)
	}
	if err := repositories.FolderRepoIns.SyncInviteeIDs(
		user, &folder, inviteeIds,
	); err != nil {
		return nil, err
	}

	res := dto.SaveFolderResponse{Folder: &folder}
	return &res, nil
}

func DeleteFolderHandler(
	c *gin.Context, user *models.User, folderId int64,
) (
	dto.ApiResponse, error,
) {
	if err := repositories.FolderRepoIns.DeleteFolder(
		user, folderId,
	); err != nil {
		return nil, err
	}

	return &dto.DeleteFolderResponse{}, nil
}
