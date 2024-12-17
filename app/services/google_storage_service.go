package services

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/innovay-software/famapp-main/app/utils"
)

var googleStorageClient *storage.Client

// Calls google storage api to do chunk upload
func GoogleStorageMultiUploader(
	objectKey string,
	localFileAbsPath string,
) error {
	ctx := context.Background()
	client, err := getGoogleStorageClient(&ctx)
	if err != nil {
		return err
	}

	bucketName := os.Getenv("GOOGLE_CLOUD_STORAGE_BUCKET_NAME")
	f, err := os.Open(localFileAbsPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucketName).Object(objectKey).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

func getGoogleStorageClient(ctx *context.Context) (*storage.Client, error) {
	if googleStorageClient != nil {
		return googleStorageClient, nil
	}
	serviceAccountJsonPath := os.Getenv("GOOGLE_SERVICE_ACCOUNT_AUTH_PATH")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", utils.GetRootAbsPath(serviceAccountJsonPath))

	client, err := storage.NewClient(*ctx)
	if err != nil {
		return nil, err
	}

	googleStorageClient = client
	return googleStorageClient, nil
}
