package jobs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/services"
	"github.com/innovay-software/famapp-main/app/utils"
)

func UploadToCloudsJob() {
	messageQueueRepo := repositories.MessageQueueRepoIns
	failedFolderFileIds := []string{}

	for {
		folderFileIdString, err := messageQueueRepo.LpopFolderFileIdFromFolderFileProcessingQueue()

		if err != nil {
			utils.LogError(err)
			break
		}

		folderFileId, err := strconv.Atoi(folderFileIdString)
		if err != nil {
			failedFolderFileIds = append(failedFolderFileIds, folderFileIdString)
			utils.LogError(err)
			continue
		}

		folderFile, err := repositories.JobRepoIns.GetFolderFile(uint64(folderFileId))
		if err != nil {
			failedFolderFileIds = append(failedFolderFileIds, folderFileIdString)
			utils.LogError(err)
			continue
		}

		if err := uploadTargetFolderFileToClouds(folderFile); err != nil {
			failedFolderFileIds = append(failedFolderFileIds, folderFileIdString)
			utils.LogError("uploadTargetFolderFileToClouds failed: ", folderFileIdString)
			utils.LogError(err)
			continue
		}

		utils.LogSuccess("uploadTargetFolderFileToClouds succeeded: ", folderFileIdString)
	}

	messageQueueRepo.RpushFailedFolderFileIdsToFailedFolderFileUploadQueue(failedFolderFileIds)
	messageQueueRepo.SetLastRunTime()
}

func uploadTargetFolderFileToClouds(folderFile *models.FolderFile) error {
	localFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.OriginalFilePath)
	if !utils.PathExists(localFileAbsPath) {
		return fmt.Errorf("unable to find originalFile for %v with path %v", folderFile.ID, folderFile.OriginalFilePath)
	}

	// Upload to HW OBS (HW stores original files only)
	if err := uploadToHwObs(folderFile); err != nil {
		utils.LogError("Upload to HW OBS failed")
		return err
	}

	// Upload to Google Drive (Google drive stores original files only)
	if err := uploadToGoogleDrive(folderFile); err != nil {
		utils.LogError("Upload to Google Drive failed")
		return err
	}

	// Compressed file
	if err := compressFolderFile(folderFile); err != nil {
		utils.LogError("Compress FolderFile failed")
		return err
	}

	// Generate thumbnail
	if err := generateFolderFileThumbnail(folderFile); err != nil {
		utils.LogError("Generate FolderFile thumbnail failed")
		return err
	}

	// Upload to Google Storage
	if err := uploadToGoogleStorage(folderFile); err != nil {
		utils.LogError("Upload to Google Storage failed")
		return err
	}

	// Delete original file
	originalFilePath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.OriginalFilePath)
	os.Remove(originalFilePath)
	folderFile.OriginalFilePath = ""
	return repositories.FolderRepoIns.SaveFolderFileModel(folderFile)
}

func uploadToHwObs(folderFile *models.FolderFile) error {
	if folderFile.HwOriginalFilePath == folderFile.OriginalFilePath {
		return nil
	}

	localFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.OriginalFilePath)
	// Upload to HW OBS
	if err := services.HwObsMultiUploader(
		folderFile.OriginalFilePath,
		localFileAbsPath,
	); err != nil {
		return err
	}
	folderFile.HwOriginalFilePath = folderFile.OriginalFilePath

	return repositories.FolderRepoIns.SaveFolderFileModel(folderFile)
}

func uploadToGoogleStorage(folderFile *models.FolderFile) error {
	if folderFile.GoogleOriginalFilePath == folderFile.OriginalFilePath &&
		folderFile.GoogleFilePath == folderFile.CompressedFilePath &&
		folderFile.GoogleThumbnailPath == folderFile.ThumbnailPath {
		return nil
	}

	localFiles := []string{
		folderFile.CompressedFilePath,
		folderFile.OriginalFilePath,
		folderFile.ThumbnailPath,
	}
	for _, item := range localFiles {
		fileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, item)
		if utils.PathExists(fileAbsPath) {
			err := services.GoogleStorageMultiUploader(
				item,
				fileAbsPath,
			)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("file not found: %v", item)
		}
	}
	folderFile.GoogleOriginalFilePath = folderFile.OriginalFilePath
	folderFile.GoogleFilePath = folderFile.CompressedFilePath
	folderFile.GoogleThumbnailPath = folderFile.ThumbnailPath

	return repositories.FolderRepoIns.SaveFolderFileModel(folderFile)
}

func uploadToGoogleDrive(folderFile *models.FolderFile) error {
	if folderFile.GoogleDriveFileID != "" {
		return nil
	}
	
	localFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.OriginalFilePath)

	fileId, err := services.GoogleDriveUploader(
		folderFile.OriginalFilePath,
		localFileAbsPath,
	)

	if err != nil {
		return err
	}
	if fileId == "" {
		return fmt.Errorf("upload to google driver returned empty fileId")
	}

	folderFile.GoogleDriveFileID = fileId
	return repositories.FolderRepoIns.SaveFolderFileModel(folderFile)
}

func compressFolderFile(folderFile *models.FolderFile) error {
	if folderFile.OriginalFilePath == "" {
		return errors.New("missing original file path")
	}
	if folderFile.CompressedFilePath == "" {
		return errors.New("missing compressed file path")
	}

	compressedFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.CompressedFilePath)

	if folderFile.FileType == "image" {
		newCompressedFileAbsPath, err := services.CompressImageFile(compressedFileAbsPath)
		if err != nil {
			return err
		}
		filename1 := filepath.Base(compressedFileAbsPath)
		filename2 := filepath.Base(newCompressedFileAbsPath)
		if filename1 != filename2 {
			folderFile.CompressedFilePath = strings.Replace(folderFile.CompressedFilePath, filename1, filename2, -1)
			if err := repositories.FolderRepoIns.SaveFolderFileModel(folderFile); err != nil {
				return err
			}
		}
		
	} else if folderFile.FileType == "video" {
		newCompressedFileAbsPath, err := services.CompressVideoFile(compressedFileAbsPath)
		if err != nil {
			return err
		}
		filename1 := filepath.Base(compressedFileAbsPath)
		filename2 := filepath.Base(newCompressedFileAbsPath)
		if filename1 != filename2 {
			folderFile.CompressedFilePath = strings.Replace(folderFile.CompressedFilePath, filename1, filename2, -1)
			if err := repositories.FolderRepoIns.SaveFolderFileModel(folderFile); err != nil {
				return err
			}
		}

	} else {
		return fmt.Errorf("invalid file type: %v", folderFile.FileType)
	}
	return nil
}

func generateFolderFileThumbnail(folderFile *models.FolderFile) error {
	if folderFile.CompressedFilePath == "" {
		return errors.New("missing compressed file path")
	}

	compressedFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.CompressedFilePath)
	ext := filepath.Ext(compressedFileAbsPath)
	thumbnailExt := ".thumbnail.jpg"

	thumbnailFileAbsPath := compressedFileAbsPath[0:len(compressedFileAbsPath) - len(ext)]
	thumbnailFileAbsPath += thumbnailExt

	if folderFile.FileType == "image" {
		err := services.GenerateImageThumbnail(compressedFileAbsPath, thumbnailFileAbsPath)
		if err != nil {
			return err
		}
	} else if folderFile.FileType == "video" {
		err := services.GenerateVideoThumbnail(compressedFileAbsPath, thumbnailFileAbsPath)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid file type: %v", folderFile.FileType)
	}

	compressedFileRelativePath := folderFile.CompressedFilePath
	thumbnailFileRelativePath := compressedFileRelativePath[0:len(compressedFileRelativePath) - len(ext)]
	thumbnailFileRelativePath += thumbnailExt

	folderFile.ThumbnailPath = thumbnailFileRelativePath
	return repositories.FolderRepoIns.SaveFolderFileModel(folderFile)
}