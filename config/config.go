package config

const (
	DefaultPort string = "8080"
	RequireJwtHeaderKey string = "Require-Jwt"
	AccessTokenKey string = "accessToken"
	RefreshTokenKey string = "refreshToken"
	AuthenticatedUserHeaderKey string = "authenticatedUser"
)

const (
	ProdStorageDir string = "storage"
	LocalStorageDir string = "storage_testing"
	UnitTestStorageDir string = "storage_unit_testing"
)

const (
	FolderFileThumbnailSize int = 250
)

var root string = ""


// Environment Overwirtes are used for different phases for production.
// Supported phases are: prod, local/dev, test(unit and integration)
var ProdServerEnvOverwrites = map[string]string{
	"STORAGE_DIR":                      ProdStorageDir,
}
var LocalServerEnvOverwrites = map[string]string{
	"DB_MAIN_DATABASE":                 "famapp_local_testing",
	"DB_READ_DATABASE":                 "famapp_local_testing",
	"MONGO_DATABASE":					"famapp_local_testing",
	"STORAGE_DIR":                      LocalStorageDir,
	"GOOGLE_CLOUD_STORAGE_BUCKET_NAME": "ijayden_testing",
	"GOOGLE_DRIVE_FOLDER_ID":           "1Ox_CrRIFsxyFfHNgYBDoOw6WYM22Ftjq",
	"HWY_BUCKET_NAME":                  "babyphotos-testing",
}
var UnitTestEnvOverwrites = map[string]string{
	"DB_MAIN_DATABASE":                 "famapp_unit_testing",
	"DB_READ_DATABASE":                 "famapp_unit_testing",
	"MONGO_DATABASE":					"famapp_unit_testing",
	"STORAGE_DIR":                      UnitTestStorageDir,
	"GOOGLE_CLOUD_STORAGE_BUCKET_NAME": "ijayden_unit_testing",
	"GOOGLE_DRIVE_FOLDER_ID":           "1fWjkNpWsXeAQtKAr_UdjhVtMmS-WcDsn",
	"HWY_BUCKET_NAME":                  "babyphotos-unit-testing",
}
