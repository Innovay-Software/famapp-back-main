/**
Everytime the job runs, it finds candidate folderFiles and uploads them to cloud

*/

package cronjob

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/innovay-software/famapp-main/app/handlers"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/services"
	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/joho/godotenv"
)

func uploadToCloud(envOverwrites map[string]string) {
	_, b, _, _ := runtime.Caller(0)
	projDir := filepath.Dir(b)
	// mainDBCon :=

	os.Setenv("projDir", projDir)
	utils.Log("projDir =", projDir)

	// Loading .env file
	utils.Log("Start API server")
	err := godotenv.Load(projDir + "/.env")
	if err != nil {
		utils.LogError("Error loading .env file", err)
	}

	// Override environment settings
	for k, v := range envOverwrites {
		utils.Log("Overriding env: ", k, "-to-", v)
		os.Setenv(k, v)
	}

	// Create storage dir if not exist
	storageAbsPath := utils.GetRootAbsPath(os.Getenv("STORAGE_DIR"))
	if !utils.PathExists(storageAbsPath) {
		os.MkdirAll(storageAbsPath, 0775)
	}

	// Repositories and handlers init
	repositories.RepoInit()
	handlers.HandlerInit()

	// Get folder files and start processing
	uploadToGoogleDrive()
}

func uploadToGoogleDrive() {

	// Set a upper limit on how many folder files this job can process at a run
	for i := 0; i < 50; i++ {
		folderFile, err := repositories.JobRepoIns.GetUploadToGoogleDriveCandidate()

		if err != nil {
			utils.LogError("Breaking uploadToGoogleDrive loop", err)
			break
		}

		localFileAbsPath := utils.GetStorageAbsPath(folderFile.Disk, folderFile.OriginalFilePath)
		if !utils.PathExists(localFileAbsPath) {
			utils.LogError("Unable to find originalFile for", folderFile.ID, "with path", folderFile.OriginalFilePath)
			continue
		}

		fileId, err := services.GoogleDriveUploader(
			folderFile.OriginalFilePath,
			localFileAbsPath,
		)

		if err != nil {
			utils.LogError("Uploading to google drive failed:", err)
			continue
		}
		if fileId == "" {
			utils.LogError("Uploading to google drive failed: no fileID")
			continue
		}

		folderFile.GoogleDriveFileID = fileId
		if err := repositories.SaveDbModel(folderFile); err != nil {
			utils.LogError("Error saving folderFile to DB")
			continue
		}

		utils.LogError(
			"", i,
			"Google Drive uploaded successful", folderFile.ID,
			folderFile.GoogleDriveFileID, folderFile.GoogleOriginalFilePath,
		)
	}
}
