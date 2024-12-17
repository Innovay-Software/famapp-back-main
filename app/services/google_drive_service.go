package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/innovay-software/famapp-main/app/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var googleDriveService *drive.Service

func GoogleDriveUploader(
	objectKey string,
	localFileAbsPath string,
) (
	string, error,
) {
	client, err := getGoogleDriveClient()
	if err != nil {
		return "", err
	}

	file, err := os.Open(localFileAbsPath)
	info, _ := file.Stat()
	if err != nil {
		return "", err
	}

	defer file.Close()

	folderId, err := getTargetFolderId(filepath.Dir(objectKey))
	if err != nil {
		return "", err
	}

	// Create File metadata
	f := &drive.File{
		Name:    info.Name(),
		Parents: []string{folderId},
	}

	// Create and upload the file
	res, err := client.Files.
		Create(f).
		Media(file).
		Do()
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func getGoogleDriveClient() (*drive.Service, error) {
	if googleDriveService != nil {
		return googleDriveService, nil
	}

	serviceAccountJsonPath := os.Getenv("GOOGLE_SERVICE_ACCOUNT_AUTH_PATH")
	serviceAccountJsonAbsPath := utils.GetRootAbsPath(serviceAccountJsonPath)

	ctx := context.Background()
	service, err := drive.NewService(
		ctx,
		option.WithCredentialsFile(serviceAccountJsonAbsPath),
		option.WithScopes(drive.DriveScope),
	)
	if err != nil {
		return nil, err
	}

	googleDriveService = service
	return service, nil
}

func getTargetFolderId(dirPath string) (string, error) {
	dirPathComponents := strings.Split(dirPath, "/")
	currentParentId := os.Getenv("GOOGLE_DRIVE_FOLDER_ID")
	if currentParentId == "" {
		return "", fmt.Errorf("missing env GOOGLE_DRIVE_FOLDER_ID")
	}

	for _, directory := range dirPathComponents {
		service, err := getGoogleDriveClient()
		if err != nil {
			return "", err
		}
		queryString := fmt.Sprintf(
			"name='%s' and mimeType = 'application/vnd.google-apps.folder' and '%s' in parents",
			directory,
			currentParentId,
		)
		res, err := service.Files.List().Q(queryString).Spaces("drive").Do()
		if err != nil {
			return "", err
		}
		if len(res.Files) > 0 {
			currentParentId = res.Files[0].Id
			// folder was already created, do nothing
		} else {
			file := drive.File{
				Name:     directory,
				MimeType: "application/vnd.google-apps.folder",
				Parents:  []string{currentParentId},
			}
			createdDir, err := service.Files.Create(&file).Do()
			if err != nil {
				return "", err
			}
			currentParentId = createdDir.Id
		}
	}
	if currentParentId == "" {
		return "", fmt.Errorf("unable to retrieve folder id")
	}
	return currentParentId, nil

}
