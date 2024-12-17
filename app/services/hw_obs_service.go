package services

import (
	"fmt"
	"log"
	"os"

	obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/innovay-software/famapp-main/app/utils"
)

const (
	multiUploadPartSize   = int64(1 * 1024 * 1024)
	multiUploadRetryLimit = 3
)

var hwObsClient *obs.ObsClient

func HwObsMultiUploader(
	objectKey string,
	localFileAbsPath string,
) error {
	utils.Log("HwObsMultiUploader:", objectKey, localFileAbsPath, multiUploadPartSize)
	obsClient, err := getHwObsClient()
	if err != nil {
		return err
	}

	var bucketName string = os.Getenv("HWY_BUCKET_NAME")
	var uploadId string = ""
	var contentLength = utils.GetFileSize(localFileAbsPath)
	var currentPartNumber int = 1
	var completedEtags []string = make([]string, 0)
	var offset int64 = 0

	// Init partial upload tasks
	inputInit := &obs.InitiateMultipartUploadInput{}
	inputInit.Bucket = bucketName
	inputInit.Key = objectKey
	outputInit, err := obsClient.InitiateMultipartUpload(inputInit)
	if err != nil {
		return err
	}
	uploadId = outputInit.UploadId

	utils.Log("UploadId:", uploadId)
	for offset < contentLength {
		utils.Log("Processing part number:", currentPartNumber)
		// Upload chunk
		inputUploadPart := &obs.UploadPartInput{}
		inputUploadPart.Bucket = bucketName
		inputUploadPart.Key = objectKey
		inputUploadPart.UploadId = uploadId
		inputUploadPart.PartNumber = currentPartNumber
		inputUploadPart.PartSize = multiUploadPartSize
		inputUploadPart.SourceFile = localFileAbsPath
		inputUploadPart.Offset = offset

		for i := 0; i < multiUploadRetryLimit+1; i++ {
			utils.Log("Upload retry:", i)
			if i == multiUploadRetryLimit {
				return fmt.Errorf(
					"upload retry limit reached for %s, partNumber=%d", objectKey, currentPartNumber,
				)
			}
			// Do the upload
			outputUploadPart, err := obsClient.UploadPart(inputUploadPart)
			if err == nil {
				utils.Log("Upload part success")
				// Record etags when completed
				completedEtags = append(completedEtags, outputUploadPart.ETag)
				break
			}
		}

		currentPartNumber += 1
		offset += multiUploadPartSize
	}

	utils.Log("Combination Stage")
	// Merge Chunks
	parts := []obs.Part{}
	for i, tag := range completedEtags {
		parts = append(parts, obs.Part{
			PartNumber: i + 1,
			ETag:       tag,
		})
	}

	inputCompleteMultipart := &obs.CompleteMultipartUploadInput{}
	inputCompleteMultipart.Bucket = bucketName
	inputCompleteMultipart.Key = objectKey
	inputCompleteMultipart.UploadId = uploadId
	inputCompleteMultipart.Parts = parts
	outputCompleteMultipart, err := obsClient.CompleteMultipartUpload(inputCompleteMultipart)
	if err != nil {
		return err
	}

	log.Printf("File uploaded, requestId:%s\n", outputCompleteMultipart.RequestId)
	log.Printf("Location:%s, Bucket:%s, Key:%s, ETag:%s\n", outputCompleteMultipart.Location, outputCompleteMultipart.Bucket, outputCompleteMultipart.Key, outputCompleteMultipart.ETag)

	return nil
}

func getHwObsClient() (*obs.ObsClient, error) {
	if hwObsClient != nil {
		return hwObsClient, nil
	}
	ak := os.Getenv("HWY_ACCESS_KEY_ID")
	sk := os.Getenv("HWY_ACCESS_KEY_SECRET")
	endPoint := os.Getenv("HWY_ENDPOINT")
	// url is babyphotos.hwobs.innovay.dev, but the prefix will be added by the library
	// endPoint = "https://hwobs.innovay.dev"

	obsClient, err := obs.New(ak, sk, endPoint)
	if err != nil {
		return nil, err
	}

	hwObsClient = obsClient
	return obsClient, err
}
