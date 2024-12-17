package cronjob

import (
	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/innovay-software/famapp-main/config"
)

// Run cron jobs
func main() {
	argsMap := utils.GetCliArgsMap()
	var envOverwrites map[string]string

	switch argsMap["-env"] {
	case "prod":
		envOverwrites = config.ProdServerEnvOverwrites
	case "local":
		envOverwrites = config.LocalServerEnvOverwrites
	default:
		envOverwrites = config.UnitTestEnvOverwrites
		utils.LogError("Missing -env cli argument:", argsMap["-env"])
	}

	uploadToCloud(envOverwrites)
}
