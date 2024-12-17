package jobs

import (
	"time"

	"github.com/innovay-software/famapp-main/app/dev"
	"github.com/innovay-software/famapp-main/app/utils"
)

func InitScheduler(initialSetupDelayInSeconds int) {
	// Wait a few seconds for system initial setup
	time.Sleep(time.Second * time.Duration(initialSetupDelayInSeconds))

	utils.Log("JobScheduler started")

	// Starts a goroutine for each scheduler
	go folderFilesCloudUploadScheduler(60)
	go systemStartupHealthCheck()
}

func folderFilesCloudUploadScheduler(intervalDelayInSeconds int) {
	for {
		didRun := foldFileCloudUploadJob(50)
		if !didRun {
			utils.Log("FolderFileCloudUploadJob run condition not satisfied, sleep until next round")
		}
		time.Sleep(time.Second * time.Duration(intervalDelayInSeconds))
	}
}

func systemStartupHealthCheck() {
	time.Sleep(time.Second * 3)
	dev.SystemHealthCheck()
}
