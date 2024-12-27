package utils

import (
	"os"
	"path"

)

// Get root absoluate path
func GetRootAbsPath(relativePath string) string {
	if relativePath != "" && relativePath[0] == '/' {
		relativePath = relativePath[1:]
	}
	return path.Join(os.Getenv("PROJECT_DIR"), relativePath)
}

// Get storage file absolute path
func GetStorageAbsPath(disk string, relativePath string) string {
	if relativePath != "" && relativePath[0] == '/' {
		relativePath = relativePath[1:]
	}
	storageDir := os.Getenv("STORAGE_DIR")
	return GetRootAbsPath(path.Join(storageDir, disk, relativePath))
}

// Construct an url to access the local file
func GetUrlPath(disk string, relativePath string) string {
	return os.Getenv("APP_HOME") + "/" + path.Join(disk, relativePath)
}
