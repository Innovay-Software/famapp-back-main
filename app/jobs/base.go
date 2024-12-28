package jobs

import (
	"time"

	"github.com/innovay-software/famapp-main/app/utils"
)

func RunJobs() {
	utils.LogSuccess("Starting UploadToCloudsJob at", time.Now())
	
	UploadToCloudsJob()
}
