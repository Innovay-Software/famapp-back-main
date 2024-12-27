package integrationTests

import (
	"encoding/base64"
	"errors"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/api"
	"github.com/innovay-software/famapp-main/app/dto"
	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/innovay-software/famapp-main/config"
	"github.com/innovay-software/famapp-main/tests"
)

// Wrapper for calling utilBase64ChunkUploadAPI, handles the complete upload.
// Returns filepath, url, error or the uploaded file
func utilBase64ChunkUploadFileWrapper(
	t *testing.T, r *gin.Engine, userToken, filepath string,
) (
	string, string, error,
) {
	f, err := os.Open(filepath)
	if err != nil {
		utils.LogError("Missing file:", filepath)
		return "", "", err
	}
	defer f.Close()
	chunkSize := 100000 // 100k Bytes
	currentOffset := int64(0)
	filename := "testing.jpeg"
	chunkedFilename := ""
	uploadedDisk := ""
	uploadedFilepath := ""
	fileUrl := ""

	for {
		_, err := f.Seek(currentOffset, 0)
		if err != nil {
			utils.LogError("Cannot seek bytes at offset:", currentOffset)
			return "", "", err
		}
		bytes := make([]byte, chunkSize) // 100k Bites
		readBytesCount, err := f.Read(bytes)
		if err != nil {
			utils.LogError("Cannot read bytes at offset:", currentOffset)
			return "", "", err
		}
		currentOffset += int64(readBytesCount)
		hasMore := readBytesCount == chunkSize && readBytesCount > 0
		readBytes := bytes[:readBytesCount]
		base64EncodedContent := base64.StdEncoding.EncodeToString(readBytes)

		res, err := utilBase64ChunkUploadAPI(r, userToken, filename, chunkedFilename, base64EncodedContent, hasMore)

		tests.AssertNil(t, err)
		if chunkedFilename == "" {
			chunkedFilename = res.ChunkedFileName
		}
		// utils.LogWarning("hasMore?", res.HasMore)
		if !res.HasMore {
			uploadedDisk = res.Document.Disk
			uploadedFilepath = res.Document.FilePath
			fileUrl = res.Document.FileUrl
			break
		}
	}

	diskFilepath := "../../" + config.UnitTestStorageDir + "/" + uploadedDisk + "/" + uploadedFilepath
	fs1, err1 := f.Stat()
	if err1 != nil {
		utils.LogError("Cannot take file1 stat")
		return "", "", err1
	}

	f2, err2 := os.Open(diskFilepath)
	if err2 != nil {
		utils.LogError("Cannot take file2 stat")
		return "", "", err2
	}

	fs2, err2 := f2.Stat()
	if err2 != nil {
		utils.LogError("Cannot take file1 stat")
		return "", "", err2
	}

	if fs1.Size() != fs2.Size() {
		return "", "", errors.New("Uploaded file size doesn't match original")
	}

	return diskFilepath, fileUrl, nil
}

// Calling API Endpoint
func utilBase64ChunkUploadAPI(
	r *gin.Engine, token, filename, chunkedFilename, base64EncodedContent string, hasMore bool,
) (
	*dto.Base64ChunkUploadFileResponse, error,
) {
	var resModel dto.Base64ChunkUploadFileResponse
	err := postRequest(
		r, "/api/v2/util/base64-chunked-upload-file",
		generateDefaultHeadersForAPIRequests(token), &api.UtilBase64ChunkUploadPathJSONRequestBody{
			FileName:             filename,
			ChunkedFileName:      chunkedFilename,
			Base64EncodedContent: base64EncodedContent,
			HasMore:              hasMore,
		},
		&resModel,
	)

	return &resModel, err
}
