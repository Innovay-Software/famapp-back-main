package dev

import (
	"fmt"
	"os"
	"strings"

	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/utils"
)

// Performs a list of system checks and returns all the errors
func SystemHealthCheck() {
	failedCount := 0
	messages := []string{
		check1_FindMissingEnvironmentVariables(&failedCount),
		check2_DbConnection(&failedCount),
		check3_ThirdPartyAPIConnection(&failedCount),
		check4_PathExistence(&failedCount),
	}

	result := "\n\nSystemHealthCheckResult:\n" + strings.Join(messages, "\n") + "\n\n"
	utils.Log(result)

	if failedCount > 0 {
		sendHealthCheckFailedReport(result)
	}
}

func sendHealthCheckFailedReport(report string) {
	// Needs Implmentation
}

func check1_FindMissingEnvironmentVariables(failedCount *int) string {
	checkName := "EnvironmentVariables"
	params := []string{
		"PROJECT_DIR",
		"GOOGLE_SERVICE_ACCOUNT_AUTH_PATH",
	}
	result := []string{}
	for _, param := range params {
		if os.Getenv(param) == "" {
			result = append(result, param)
		}
	}

	if len(result) > 0 {
		(*failedCount)++
		return errorMessage(
			checkName,
			fmt.Errorf("missing: "+strings.Join(result, ",")),
		)
	}
	return successMessage(checkName)
}

func check2_DbConnection(failedCount *int) string {
	checkName := "DatabaseConnection"
	users, err := repositories.UserRepoIns.FindAllUser()
	if err != nil {
		return errorMessage(checkName, fmt.Errorf("error when querying database: %s", err.Error()))
	}
	utils.Log("Found ", len(*users), "users", (*users)[0].Mobile)
	configCount := repositories.ConfigRepoIns.CountConfigRecords()
	if configCount <= 0 {
		(*failedCount)++
		return errorMessage(checkName, fmt.Errorf("error when querying database, record not found"))
	}
	return successMessage(checkName)
}

func check3_ThirdPartyAPIConnection(failedCount *int) string {
	checkName := "ThirdPartyAPIConnection_NeedsImplementation"
	if checkName == "test" {
		(*failedCount)++
	}
	return successMessage(checkName)
}

func check4_PathExistence(failedCount *int) string {
	checkName := "PathExistence"
	storageDir := os.Getenv("STORAGE_DIR")
	if !utils.PathExists(utils.GetRootAbsPath(storageDir)) {
		(*failedCount)++
		return errorMessage(checkName, fmt.Errorf("storage path not exist: %s", utils.GetRootAbsPath(storageDir)))
	}
	return successMessage(checkName)
}

func successMessage(name string) string {
	return "***  SystemCheck." + name + " succeeded"
}

func errorMessage(name string, err error) string {
	return "***  !!!SystemCheck." + name + " failed: " + err.Error()
}
