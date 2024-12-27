package folderFile

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/utils"
)

type FolderFileMediaType int

const (
	FolderFileMediaCompressedFile FolderFileMediaType = iota
	FolderFileMediaOriginalFile
	FolderFileMediaThumbnailFile
)

// Display folderFile.FilePath image file
func DisplayFolderFileCompressedHandler(
	c *gin.Context, user *models.User, folderFileId uint64,
) (
	string, error,
) {
	return getFolderFileDisplayFilepath(c, FolderFileMediaCompressedFile, user, folderFileId)
}

// Display folderFile.ThumbnailPath image file
func DisplayFolderFileThumbnailHandler(
	c *gin.Context, user *models.User, folderFileId uint64,
) (
	string, error,
) {
	filepath, err := getFolderFileDisplayFilepath(c, FolderFileMediaThumbnailFile, user, folderFileId)
	if err == nil {
		return filepath, nil
	}
	// If thumbnail is not available, return the compressed file
	return getFolderFileDisplayFilepath(c, FolderFileMediaCompressedFile, user, folderFileId)
}

// Display folderFile.Originalpath image file
func DisplayFolderFileOriginalHandler(
	c *gin.Context, user *models.User, folderFileId uint64,
) (
	string, error,
) {
	filepath, err := getFolderFileDisplayFilepath(c, FolderFileMediaOriginalFile, user, folderFileId)
	if err == nil {
		return filepath, nil
	}
	// If original file is not available, return the compressed file
	return getFolderFileDisplayFilepath(c, FolderFileMediaCompressedFile, user, folderFileId)
}

func getFolderFileDisplayFilepath(
	c *gin.Context, target FolderFileMediaType, user *models.User, folderFileId uint64,
) (
	string, error,
) {
	if folderFile, err := repositories.FolderRepoIns.GetFolderFileByFieldName(
		user, "id", strconv.Itoa(int(folderFileId)),
	); err == nil {
		relativePath := ""
		switch target {
		case FolderFileMediaCompressedFile:
			relativePath = folderFile.CompressedFilePath
		case FolderFileMediaThumbnailFile:
			relativePath = folderFile.ThumbnailPath
		case FolderFileMediaOriginalFile:
			relativePath = folderFile.OriginalFilePath
		}

		if relativePath != "" {
			filepath := utils.GetStorageAbsPath(folderFile.Disk, relativePath)
			if utils.PathExists(filepath) {
				// c.File(filepath)
				return filepath, nil
			}
		} else {
			return "", fmt.Errorf("invalid target type: %d", target)
		}
	}

	return "", errors.ApiErrorPermissionDenied
}
