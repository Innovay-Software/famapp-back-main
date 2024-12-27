package folderFile

import (
	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
	apiErrors "github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
)

func GetFolderFilesBeforeTakenOnHandler(
	c *gin.Context, user *models.User, folderId uint64, pivotDate string, beforeMicroTimestamp int64,
) (
	dto.ApiResponse, error,
) {
	folderFiles, folder, hasMore, err := repositories.FolderRepoIns.GetFolderFilesBeforeTakenOn(
		user, folderId, pivotDate, beforeMicroTimestamp,
	)
	if err != nil {
		return nil, err
	}

	res := dto.GetFolderFilesBeforeTakenOnResponse{
		FolderFiles: folderFiles,
		Folder:      folder,
		HasMore:     hasMore,
	}

	return &res, nil
}

func GetFolderFilesAfterTakenOnHandler(
	c *gin.Context, user *models.User, folderId uint64, pivotDate string, afterMicroTimestamp int64,
) (
	dto.ApiResponse, error,
) {
	folderFiles, folder, hasMore, err := repositories.FolderRepoIns.GetFolderFilesAfterTakenOn(
		user, folderId, pivotDate, afterMicroTimestamp,
	)
	if err != nil {
		return nil, err
	}

	res := dto.GetFolderFilesAfterTakenOnResponse{
		FolderFiles: folderFiles,
		Folder:      folder,
		HasMore:     hasMore,
	}

	return &res, nil
}

func UpdateSingleFolderFileHandler(
	c *gin.Context, user *models.User, folderFileId uint64, isPrivate *bool, remark *string,
) (
	dto.ApiResponse, error,
) {
	var ff models.FolderFile
	if err := repositories.QueryDbModelByPrimaryId(&ff, folderFileId); err != nil {
		return nil, err
	}
	if !user.IsAdmin() && ff.OwnerID != user.ID {
		var f models.Folder
		if err := repositories.QueryDbModelByPrimaryId(&f, ff.FolderID); err != nil {
			return nil, err
		}
		if f.OwnerID != user.ID {
			return nil, apiErrors.ApiErrorPermissionDenied
		}
	}

	if isPrivate != nil {
		ff.IsPrivate = *isPrivate
	}
	if remark != nil {
		ff.Remark = *remark
	}

	repositories.SaveDbModel(&ff)
	return &dto.UpdateSingleFolderFileResponse{}, nil
}

func UpdateMultipleFolderFilesHandler(
	c *gin.Context, user *models.User, folderFileIds []uint64, newFolderId *uint64, newTakenOnTimeStamp *int64,
) (
	dto.ApiResponse, error,
) {
	if newFolderId != nil {
		// Move folder files
		err := repositories.FolderRepoIns.MoveFolderFiles(user, &folderFileIds, *newFolderId, 100)
		if err != nil {
			return nil, err
		}
	}
	if newTakenOnTimeStamp != nil {
		// Set timestamp
		err := repositories.FolderRepoIns.RescheduleFolderFiles(user, &folderFileIds, *newTakenOnTimeStamp, 100)
		if err != nil {
			return nil, err
		}
	}

	return &dto.UpdateMultipleFolderFileResponse{}, nil
}

func DeleteFolderFilesHandler(
	c *gin.Context, user *models.User, folderId uint64, folderFileIds []uint64,
) (
	dto.ApiResponse, error,
) {

	failedIds := []uint64{}
	// needs to delete each one individually to check for user permision
	for _, item := range folderFileIds {
		if err := repositories.FolderRepoIns.DeleteFolderFile(user, item); err != nil {
			failedIds = append(failedIds, item)
		}
	}

	if len(failedIds) > 0 {
		return nil, apiErrors.ApiErrorPermissionDenied
	}

	return &dto.DeleteFolderFileResponse{}, nil
}
