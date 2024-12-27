package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/innovay-software/famapp-main/app/api"
	"github.com/innovay-software/famapp-main/app/handlers"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/config"

	// "github.com/innovay-software/famapp-main/app/jobs"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/joho/godotenv"
)

// A struct to hold the current app
type ApiApp struct {
	// Is Production build
	IsProd bool
	// Is Debug mode
	IsDebug bool
	// Scheduler initial start delay in seconds
	SchedulerInitialDelay int
	// Environment variable overwrites
	EnvOverwrites map[string]string
}

// Production Server
func InitApiProdServer(projDir string) (*gin.Engine, string) {
	productionApp := ApiApp{
		IsProd:                true,
		IsDebug:               true,
		SchedulerInitialDelay: 5,
		EnvOverwrites:         config.ProdServerEnvOverwrites,
	}
	return productionApp.initServer(projDir)
}

// Local Server
func InitApiLocalServer(projDir string) (*gin.Engine, string) {
	localApp := ApiApp{
		IsProd:                false,
		IsDebug:               true,
		SchedulerInitialDelay: -1, // -1 to disable job scheduler
		EnvOverwrites:         config.LocalServerEnvOverwrites,
	}
	return localApp.initServer(projDir)
}

// Unit and Integration Test Server
func InitApiIntegrationTestServer(projDir string) (*gin.Engine, string) {
	integrationTestApp := ApiApp{
		IsProd:                false,
		IsDebug:               true,
		SchedulerInitialDelay: -1, // -1 to disable job scheduler
		EnvOverwrites:         config.UnitTestEnvOverwrites,
	}
	return integrationTestApp.initServer(projDir)
}

// Start gin server and returns it with a port number
func (a ApiApp) initServer(projDir string) (*gin.Engine, string) {

	// Load dotEnv and other env variables into os.Setenv
	a.loadEnvironmentVariables(projDir)
	utils.LogSuccess("Using projDir =", os.Getenv("PROJECT_DIR"))

	// Get target port string
	port := getPortString()
	utils.LogSuccess("Using port =", port)

	// Create storage dir if not exist
	storageDir := os.Getenv("STORAGE_DIR")
	a.createStaticAndStorageDirsIfNotExist(storageDir)

	// Repositories and handlers initialization
	repositories.RepoInit()
	handlers.HandlerInit()

	// Check if default admin user exists:
	a.createDefaultAdminUser()

	// Get gin router
	r := gin.Default()

	// Static files
	staticMap := map[string]string{
		"static":      utils.GetRootAbsPath("static"),
		"user-upload": utils.GetRootAbsPath(storageDir + "/user-upload"),
		"avatars":     utils.GetRootAbsPath(storageDir + "/avatars"),
		"album-cover": utils.GetRootAbsPath(storageDir + "/album-cover"),
	}
	for k, v := range staticMap {
		r.Static(k, v)
	}

	api.RegisterRoutes(r, a.IsProd)
	return r, port
}

// Load all variables and overwrites into os.Setenv
func (a ApiApp) loadEnvironmentVariables(projDir string) {
	err := godotenv.Load(projDir + "/.env")
	if err != nil {
		utils.LogError("Error loading .env file", err)
	}

	// Override environment settings
	for k, v := range a.EnvOverwrites {
		utils.LogSuccess("Overriding env: ", k, "-to-", v)
		os.Setenv(k, v)
	}

	os.Setenv("PROJECT_DIR", projDir)
	if a.IsProd {
		os.Setenv("IS_PROD", "1")
	} else {
		os.Setenv("IS_PROD", "0")
	}
	if a.IsDebug {
		os.Setenv("IS_DEBUG", "1")
	} else {
		os.Setenv("IS_DEBUG", "0")
	}
}

// Check if the necessary dirs exist, if not, create them
func (a ApiApp) createStaticAndStorageDirsIfNotExist(storageDir string) {
	storageAbsPath := utils.GetRootAbsPath(storageDir)
	if !utils.PathExists(storageAbsPath) {
		os.MkdirAll(storageAbsPath, 0775)
	}
}

// Gets the port from env, if not found, use the default from config
func getPortString() string {
	port := os.Getenv("PORT")
	if port == "" {
		utils.LogWarning("Missing port, defaulting to", config.DefaultPort)
		port = config.DefaultPort
	}
	return port
}

// Create default admin user
func (a ApiApp) createDefaultAdminUser() {
	_, err := repositories.UserRepoIns.FindUserByField("id", "1")
	if err != nil {
		// Not able to find admin user
		uuid, _ := uuid.FromBytes([]byte(os.Getenv("ADMIN_UUID")))
		user := &models.User{
			BaseDbModel: models.BaseDbModel{
				ID: 1,
			},
			UUID:           uuid,
			Role:           "admin",
			Name:           os.Getenv("ADMIN_NAME"),
			Email:          os.Getenv("ADMIN_EMAIL"),
			Mobile:         os.Getenv("ADMIN_MOBILE"),
			LockerPasscode: os.Getenv("ADMIN_LOCKER_PIN"),
			Avatar:         "",
		}
		user.SetPassword(os.Getenv("ADMIN_PASSWORD"))
		repositories.UserRepoIns.CreateUser(user)
	}
}
