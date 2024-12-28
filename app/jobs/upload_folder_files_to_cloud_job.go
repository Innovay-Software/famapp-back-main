package jobs

// import (
// 	"fmt"
// 	"strings"
// 	"time"

// 	"github.com/innovay-software/famapp-main/app/models"
// 	"github.com/innovay-software/famapp-main/app/repositories"
// 	"github.com/innovay-software/famapp-main/app/services"
// 	"github.com/innovay-software/famapp-main/app/utils"
// )

// // Find target folder files and try to upload to cloud storage for backup and CDN
// func foldFileCloudUploadJob(candidateLimit int) bool {
// 	if !canStartCloudUploadProcess(1 * 60) {
// 		return false
// 	}

// 	if err := startHwObsCloudUploadProcess(candidateLimit); err != nil {
// 		utils.Log("Error uploading to HW OBS:", err)
// 	}
// 	if err := startGoogleDriveCloudUploadProcess(candidateLimit); err != nil {
// 		utils.Log("Error uploading to Google Drive:", err)
// 	}
// 	if err := startGoogleStorageCloudUploadProcess(candidateLimit); err != nil {
// 		utils.Log("Error uploading to Google Storage:", err)
// 	}

// 	return true
// }

// // Checks if the cloud upload process can be triggered right now
// func canStartCloudUploadProcess(minimumDelayInSeconds int64) bool {
// 	nowUnix := time.Now().UTC().Unix()
// 	lastUploadTimeUnix := repositories.ConfigRepoIns.GetUploadTime().Unix()

// 	return nowUnix-lastUploadTimeUnix > minimumDelayInSeconds
// }

// // Starts the HW OBS cloud upload process
// func startHwObsCloudUploadProcess(candidateLimit int) error {
// 	repo := repositories.CloudUploadRepoIns

// 	// Get the first 50 folder files to upload to hw obs
// 	folderFiles, err := repo.GetHwObsCandidates(candidateLimit)
// 	if err != nil {
// 		return err
// 	}

// 	if err := updateFolderFilesToProcessingStatus("hw_original_file_path", folderFiles); err != nil {
// 		return err
// 	}

// 	utils.Log("Cloud Upload: found", len(*folderFiles), "files to upload to HW OBS")

// 	for _, folderFile := range *folderFiles {
// 		utils.Log("**HW OBS Processing", folderFile.ID)
// 		localFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.CompressedFilePath)
// 		if !utils.PathExists(localFileAbsPath) {
// 			utils.Log("File not exist:", localFileAbsPath)
// 			continue
// 		}
// 		var err error
// 		err = services.HwObsMultiUploader(
// 			folderFile.OriginalFilePath,
// 			localFileAbsPath,
// 		)
// 		if err == nil {
// 			folderFile.HwOriginalFilePath = folderFile.OriginalFilePath
// 			err = repositories.FolderRepoIns.SaveFolderFileModel(&folderFile)
// 		}
// 		if err != nil {
// 			utils.Log("Uploading folder file", folderFile.ID, "error: ", err)
// 		}
// 	}

// 	return nil
// }

// // Starts the Google Storage cloud upload process
// func startGoogleStorageCloudUploadProcess(candidateLimit int) error {
// 	repo := repositories.CloudUploadRepoIns

// 	// Get the first 50 folder files to upload to hw obs
// 	folderFiles, err := repo.GetGoogleStorageCandidates(candidateLimit)
// 	if err != nil {
// 		utils.Log(err)
// 		return err
// 	}

// 	if err := updateFolderFilesToProcessingStatus("google_file_path", folderFiles); err != nil {
// 		utils.Log(err)
// 		return err
// 	}

// 	utils.Log("Cloud Upload: found", len(*folderFiles), "files to upload to Google Storage")

// 	for _, folderFile := range *folderFiles {
// 		utils.Log("**Google Storage Processing", folderFile.ID)

// 		// Refine folder files metadatas
// 		if err := refineFolderFileMetadata(&folderFile); err != nil {
// 			utils.Log("err:", err)
// 		}

// 		// Generate thumbnail
// 		if err := generateFolderFileThumbnail(&folderFile); err != nil {
// 			utils.Log("err:", err)
// 		}

// 		// Compress folder file
// 		if err := compressFolderFile(&folderFile); err != nil {
// 			utils.Log("err:", err)
// 		}

// 		// Upload to Google
// 		successful := true
// 		localFiles := []string{
// 			folderFile.CompressedFilePath,
// 			folderFile.OriginalFilePath,
// 			folderFile.ThumbnailPath,
// 		}
// 		for _, item := range localFiles {
// 			fileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, item)
// 			if utils.PathExists(fileAbsPath) {
// 				utils.Log("Uploading", item, "to google cloud storage")
// 				err := services.GoogleStorageMultiUploader(
// 					item,
// 					fileAbsPath,
// 				)
// 				if err != nil {
// 					successful = false
// 				}
// 			} else {
// 				utils.Log("File not found:", item)
// 			}
// 		}
// 		if successful {
// 			folderFile.GoogleOriginalFilePath = folderFile.OriginalFilePath
// 			folderFile.GoogleFilePath = folderFile.CompressedFilePath
// 			folderFile.GoogleThumbnailPath = folderFile.ThumbnailPath
// 			// Save folder file to database
// 			if err := repositories.FolderRepoIns.SaveFolderFileModel(&folderFile); err != nil {
// 				utils.Log("Uploading folder file", folderFile.ID, "error: ", err)
// 			}
// 		}
// 	}

// 	return nil
// }

// // Starts the Google Drive cloud upload process
// func startGoogleDriveCloudUploadProcess(candidateLimit int) error {
// 	repo := repositories.CloudUploadRepoIns

// 	// Get the first 50 folder files to upload to hw obs
// 	folderFiles, err := repo.GetGoogleDriveCandidates(candidateLimit)
// 	if err != nil {
// 		utils.Log(err)
// 		return err
// 	}

// 	if err := updateFolderFilesToProcessingStatus("google_drive_file_id", folderFiles); err != nil {
// 		utils.Log(err)
// 		return err
// 	}

// 	utils.Log("Cloud Upload: found", len(*folderFiles), "files to upload to Google Drive")

// 	for _, folderFile := range *folderFiles {
// 		utils.Log("**Google Drive Processing", folderFile.ID, "google original: ", folderFile.GoogleOriginalFilePath)
// 		localFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.OriginalFilePath)
// 		if !utils.PathExists(localFileAbsPath) {
// 			continue
// 		}

// 		fileId, err := services.GoogleDriveUploader(
// 			folderFile.OriginalFilePath,
// 			localFileAbsPath,
// 		)
// 		if err != nil {
// 			utils.Log("Uploading to google drive failed:", err)
// 			continue
// 		}
// 		if fileId == "" {
// 			utils.Log("Uploading to google drive failed: no fileID")
// 			continue
// 		}
// 		folderFile.GoogleDriveFileID = fileId
// 		if err := repositories.FolderRepoIns.SaveFolderFileModel(&folderFile); err != nil {
// 			continue
// 		}
// 		utils.Log(
// 			"Google Drive uploaded successful", folderFile.ID,
// 			folderFile.GoogleDriveFileID, folderFile.GoogleOriginalFilePath,
// 		)
// 	}

// 	return nil
// }

// func updateFolderFilesToProcessingStatus(fieldName string, folderFiles *[]models.FolderFile) error {
// 	repo := repositories.CloudUploadRepoIns
// 	ids := []int{}
// 	for _, folderFile := range *folderFiles {
// 		ids = append(ids, int(folderFile.ID))
// 	}
// 	if err := repo.BatchUploadFolderFile(fieldName, "processing", ids); err != nil {
// 		utils.Log(err)
// 		return err
// 	}
// 	return nil
// }

// func refineFolderFileMetadata(folderFile *models.FolderFile) error {
// 	utils.Log("refineFolderFileMetadata")
// 	targetFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.OriginalFilePath)
// 	if !utils.PathExists(targetFileAbsPath) {
// 		displayingFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.CompressedFilePath)
// 		if !utils.PathExists(displayingFileAbsPath) {
// 			return fmt.Errorf(
// 				"original and display file not exists: %s, %s",
// 				folderFile.OriginalFilePath, folderFile.CompressedFilePath,
// 			)
// 		}
// 		targetFileAbsPath = displayingFileAbsPath
// 	}
// 	if folderFile.IsImage() {
// 		width, height := utils.GetImageDimeision(targetFileAbsPath)
// 		folderFile.Metadata["dimension"] = fmt.Sprint(width, "x", height)
// 	}
// 	if folderFile.IsVideo() {
// 		firstFrameAbsPath := utils.ChangeFileExtension(targetFileAbsPath, "jpg")
// 		_, err := utils.ExtractVideoFirstFrameAsJpg(targetFileAbsPath, firstFrameAbsPath)
// 		if err != nil {
// 			return err
// 		}
// 		if utils.PathExists(firstFrameAbsPath) {
// 			width, height := utils.GetImageDimeision(firstFrameAbsPath)
// 			folderFile.Metadata["dimension"] = fmt.Sprint(width, "x", height)
// 		}
// 	}
// 	return nil
// }

// func compressFolderFile(folderFile *models.FolderFile) error {
// 	utils.Log("compressFolderFile")
// 	originalFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.OriginalFilePath)
// 	originalExists := utils.PathExists(originalFileAbsPath)

// 	targetFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.CompressedFilePath)
// 	targetExists := utils.PathExists(targetFileAbsPath)

// 	if !originalExists && !targetExists {
// 		return fmt.Errorf("missing target file: %s, %s", originalFileAbsPath, targetFileAbsPath)
// 	}
// 	if originalExists && !targetExists {
// 		utils.DuplicateFile(originalFileAbsPath, targetFileAbsPath)
// 	} else if !originalExists && targetExists {
// 		utils.DuplicateFile(targetFileAbsPath, originalFileAbsPath)
// 	}

// 	compressedFileAbsPath := ""
// 	if folderFile.IsImage() {
// 		filePath, err := utils.CompressImageToJpgWithMaxSize(
// 			targetFileAbsPath, 1080,
// 		)
// 		if err != nil {
// 			return err
// 		}
// 		compressedFileAbsPath = filePath
// 		syncImageExifData(compressedFileAbsPath, folderFile)
// 	}
// 	if folderFile.IsVideo() {
// 		// dimension := folderFile.Metadata["dimension"].(string)
// 		// dimensionParts := strings.Split(dimension, "x")
// 		// if len(dimensionParts) < 2 {
// 		// 	return fmt.Errorf("unable to extract width and height from %s", dimension)
// 		// }

// 		// width, err := strconv.Atoi(dimensionParts[0])
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		// height, err := strconv.Atoi(dimensionParts[1])
// 		// if err != nil {
// 		// 	return err
// 		// }

// 		filePath, err := utils.CompressVideoToMp4FullHD(targetFileAbsPath)
// 		if err != nil {
// 			return err
// 		}
// 		compressedFileAbsPath = filePath
// 	}

// 	tempParts := strings.Split(compressedFileAbsPath, "/"+folderFile.Disk+"/")
// 	if len(tempParts) > 1 {
// 		folderFile.CompressedFilePath = tempParts[len(tempParts)-1]
// 	} else if len(tempParts) == 1 {
// 		folderFile.CompressedFilePath = compressedFileAbsPath
// 	}

// 	return nil
// }

// // Generate thumbnail file
// func generateFolderFileThumbnail(folderFile *models.FolderFile) error {
// 	utils.Log("generateFolderFileThumbnail")

// 	// Generate thumbnail from compressFile or originalFile, which ever exists first
// 	filePaths := []string{folderFile.CompressedFilePath, folderFile.OriginalFilePath}
// 	targetFileAbsPath := ""
// 	for _, item := range filePaths {
// 		absPath := utils.GetStorageAbsPath(folderFile.Disk, item)
// 		if utils.PathExists(absPath) {
// 			targetFileAbsPath = absPath
// 			break
// 		}
// 	}

// 	// if could not find a reference image to generate thumbnail from, return an error
// 	if targetFileAbsPath == "" {
// 		return fmt.Errorf("missing target file: %s", strings.Join(filePaths, ", "))
// 	}

// 	// Get thumbnail path and check if a thumbnail already exists
// 	thumbnailRelativePath := folderFile.GenerateThumbnailPath()
// 	thumbnailAbsPath := utils.GetStorageAbsPath(folderFile.Disk, thumbnailRelativePath)
// 	if utils.PathExists(thumbnailAbsPath) {
// 		// If thumbnail already exists, do not regenerate
// 		folderFile.ThumbnailPath = thumbnailRelativePath
// 		return nil
// 	}

// 	if folderFile.IsImage() {
// 		// Generate thumbnail
// 		if _, err := utils.GenerateThumbnailJpg(targetFileAbsPath, thumbnailAbsPath); err != nil {
// 			return err
// 		}
// 		// Sync EXIF data
// 		if err := syncImageExifData(thumbnailAbsPath, folderFile); err != nil {
// 			return err
// 		}
// 	}
// 	if folderFile.IsVideo() {
// 		targetFileAbsPath = utils.GetStorageAbsPath(folderFile.Disk, folderFile.OriginalFilePath)
// 		firstFrameAbsPath := utils.ChangeFileExtension(targetFileAbsPath, "jpg")
// 		if _, err := utils.GenerateThumbnailJpg(firstFrameAbsPath, thumbnailAbsPath); err != nil {
// 			return err
// 		}
// 	}

// 	// Save to repository
// 	folderFile.ThumbnailPath = thumbnailRelativePath
// 	return repositories.FolderRepoIns.SaveFolderFileModel(folderFile)
// }

// // Extract some key EXIF data from the oroginal file and apply them to the target file
// // Key EXIF data include: Orientation
// func syncImageExifData(targetFileAbsPath string, folderFile *models.FolderFile) error {
// 	orientationExifString := ""
// 	if exifVal, exists := (*folderFile).Metadata["exif"]; exists {
// 		exifMap := exifVal.(map[string]any)
// 		if orientationValue, exists := exifMap["Orientation"]; exists {
// 			if orientationString, ok := orientationValue.(string); ok {
// 				orientationExifString = orientationString
// 			}
// 		}
// 	}
// 	return utils.SetImageExif(targetFileAbsPath, orientationExifString)
// }
