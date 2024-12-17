package folderFile

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
	"github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/services"
	"github.com/innovay-software/famapp-main/app/utils"
)

var chunkUploadDisk = "chunk-upload"

func FolderFileChunkUploadInitUploadIdHandler(
	c *gin.Context, user *models.User,
) (
	dto.ApiResponse, error,
) {
	uploadInstance := models.Upload{
		UserID:   user.ID,
		Disk:     chunkUploadDisk,
		FileType: "others",
		ShotAt:   time.Now(),
	}
	if err := repositories.SaveDbModel(&uploadInstance); err != nil {
		return nil, err
	}

	res := dto.FolderFileGetChunkUploadFileIdResponse{
		UploadId: uploadInstance.ID,
	}

	return &res, nil
}

func FolderFileChunkUploadHandler(
	c *gin.Context, user *models.User, folderId int64, uploadId string,
	hasmore bool, filename string, chunkindex int,
) (
	dto.ApiResponse, error,
) {
	// Get folder
	var folder models.Folder
	if err := repositories.QueryDbModelByPrimaryId(&folder, folderId); err != nil {
		return nil, errors.ApiErrorPermissionDenied
	}

	// Check for post permission
	if !repositories.FolderRepoIns.HasFolderPostPermission(user, &folder) {
		return nil, errors.ApiErrorPermissionDenied
	}

	// Get filename, chunkindex, and hasmore headers
	if filename == "" {
		return nil, errors.ApiErrorParamInvalid
	}

	// Get rawdata
	byteList, err := c.GetRawData()
	if err != nil {
		return nil, errors.ApiErrorParamInvalid
	}

	// Write/Append to file
	chunkFileAbsPath := utils.GetStorageAbsPath(chunkUploadDisk, uploadId+"_"+filename)

	f, err := os.OpenFile(chunkFileAbsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	if chunkindex == 1 {
		if err := f.Truncate(0); err != nil {
			return nil, err
		}
		if _, err := f.Seek(0, 0); err != nil {
			return nil, err
		}
	}

	if _, err := f.Write(byteList); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	// Update upload time
	repositories.ConfigRepoIns.UpdateUploadTime()
	if hasmore {
		return &dto.FolderFileChunkUploadFileResponse{}, nil
	}

	// else if no more chunks are expected, delete upload record and process the uploaded file
	var uploadInstance models.Upload
	if uploadIdInt, err := strconv.Atoi(uploadId); err == nil {
		if err := repositories.QueryDbModelByPrimaryId(&uploadInstance, int64(uploadIdInt)); err == nil {
			repositories.DeleteDbModel(uploadInstance)
		}
	}

	if err := processUploadedChunkFileAsFolderFile(
		user, &folder, chunkFileAbsPath,
	); err != nil {
		return nil, err
	}

	return &dto.FolderFileChunkUploadFileResponse{}, nil
}

// Processes uploaded chunk file and save to database as a folder file record
func processUploadedChunkFileAsFolderFile(
	user *models.User, folder *models.Folder, chunkFileAbsPath string,
) error {
	preventDuplicateFiles := false

	// Get file MD5 value
	fileMd5, err := utils.GenerateFileMd5(chunkFileAbsPath)
	if err != nil {
		return errors.ApiErrorSystem
	}

	// Check if there are other files with same md5 value
	if preventDuplicateFiles {
		if _, err := repositories.FolderRepoIns.GetActiveFolderFileWithMd5(
			int64(folder.ID), fileMd5,
		); err == nil {
			utils.DeleteFile(chunkFileAbsPath)
			return errors.ApiError{Code: -1, Message: "Duplicate File"}
		}
	}

	// Get Metadata data
	ext := filepath.Ext(chunkFileAbsPath)
	filetype := utils.FileExtToFileType(ext)
	metadataMap := services.ExtractFileMetadata(chunkFileAbsPath)
	targetDatetime := (*metadataMap)["shot_at_date_time"].(time.Time)

	// Update target date time micro seconds
	targetUnixSecondTimestamp := targetDatetime.UTC().Unix()
	fileCountAtTargetSecond := repositories.FolderRepoIns.CountFilesForTargetUnixTimestamp(targetUnixSecondTimestamp)
	targetDatetime = time.Unix(targetUnixSecondTimestamp, (fileCountAtTargetSecond+1)*1000)
	(*metadataMap)["shot_at_date_time"] = targetDatetime

	// Get save path
	albumDisk := "album-general"
	datePath := strings.ReplaceAll(targetDatetime.Format(time.DateOnly), "-", "/")
	dateString := strings.ReplaceAll(targetDatetime.Format(time.DateOnly), "-", "_")
	relativeDirPath := fmt.Sprintf("album%d/%d/%s", folder.ID, user.ID, datePath)
	absoluteDirPath := utils.GetStorageAbsPath(albumDisk, relativeDirPath)
	if !utils.PathExists(absoluteDirPath) {
		os.MkdirAll(absoluteDirPath, 0775)
	}
	tempString := strconv.Itoa(rand.Intn(9999-1001) + 1001)
	filename := dateString + "_" + tempString + ext
	originalFilename := dateString + "_" + tempString + ".original" + ext

	fileAbsPath := absoluteDirPath + "/" + filename
	originalFileAbsPath := absoluteDirPath + "/" + originalFilename

	// Move chunk upload file to its target destination path
	if err := os.Rename(chunkFileAbsPath, fileAbsPath); err != nil {
		return errors.ApiErrorSystem
	}

	// Make a copy for original file
	if err := utils.DuplicateFile(fileAbsPath, originalFileAbsPath); err != nil {
		return errors.ApiErrorSystem
	}

	// Create folder file model instance
	dateTimeNow := time.Now().Unix()
	folderFile := models.FolderFile{
		OwnerID:  user.ID,
		FolderID: folder.ID,
		Disk:     albumDisk,
		FileName: fileMd5 + "." + strconv.Itoa(int(dateTimeNow)) + ext,
		FileType: filetype,

		OriginalFilePath:   strings.Split(originalFileAbsPath, "/"+albumDisk+"/")[1],
		CompressedFilePath: strings.Split(fileAbsPath, "/"+albumDisk+"/")[1],
		ThumbnailPath:      "",

		Remark:   "",
		Metadata: *metadataMap,
		ShotAt:   targetDatetime,
	}

	if err := repositories.FolderRepoIns.CreateFolderFile(&folderFile); err != nil {
		return errors.ApiErrorSystem
	}
	return nil
}
