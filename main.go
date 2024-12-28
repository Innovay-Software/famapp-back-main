package main

import (
	"path/filepath"
	"runtime"

	"github.com/innovay-software/famapp-main/app"
	"github.com/innovay-software/famapp-main/app/jobs"
	"github.com/innovay-software/famapp-main/app/utils"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	projDir := filepath.Dir(b)
	argsMap := utils.GetCliArgsMap()

	switch argsMap["-role"] {
	case "webserver":
		utils.LogSuccess("starting webserver")
		webserver(&argsMap, projDir)

	case "cronjob":
		utils.LogSuccess("starting cronjob")
		cronjob(&argsMap, projDir)
		
	default:
		utils.LogError("Missing -role cli argument:", argsMap["-role"])
	}
}

func webserver(argsMapPointer *map[string]string, projDir string) {
	argsMap := *argsMapPointer
	switch argsMap["-env"] {
	case "prod":
		utils.LogSuccess("Starting Prod Server:", projDir)
		ginInstance, port := app.InitApiProdServer(projDir)
		ginInstance.Run(":" + port)
	case "local":
		utils.LogSuccess("Start Local Server:", projDir)
		ginInstance, port := app.InitApiLocalServer(projDir)
		utils.PrintAllRoutes(ginInstance)
		ginInstance.Run(":" + port)
	default:
		utils.LogError("Missing -env cli argument:", argsMap["-env"])
	}

	utils.LogError("Gin server terminated...")
}

func cronjob(argsMapPointer *map[string]string, projDir string) {
	argsMap := *argsMapPointer
	switch argsMap["-env"] {
	case "prod":
		utils.LogSuccess("Starting Prod Server:", projDir)
		app.InitApiProdServer(projDir)
		jobs.RunJobs()

	case "local":
		utils.LogSuccess("Start Local Server:", projDir)
		app.InitApiLocalServer(projDir)
		jobs.RunJobs()
		
	default:
		utils.LogError("Missing -env cli argument:", argsMap["-env"])
	}
}

