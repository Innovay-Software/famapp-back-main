package integrationTests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/api"
	"github.com/innovay-software/famapp-main/app/dto"
)

// Calling API Endpoint
func folderFileChunkUploadInitUploadIdAPI(
	gin *gin.Engine, token string,
) (*dto.FolderFileGetChunkUploadFileIdResponse, error) {
	var resModel dto.FolderFileGetChunkUploadFileIdResponse
	err := postRequest(
		gin, "/api/v2/folder-files/chunk-upload-folder-file-init-upload-id",
		generateDefaultHeadersForAPIRequests(token),
		nil, &resModel,
	)
	return &resModel, err
}

// Calling API Endpoint
func uploadFileToFolderFile(
	gin *gin.Engine, accessToken string, folderID int64, filePath string,
) (
	*dto.FolderFileChunkUploadFileResponse, error,
) {

	uploadId := int64(0)
	{
		res, err := folderFileChunkUploadInitUploadIdAPI(gin, accessToken)
		if err != nil {
			return nil, err
		}
		uploadId = res.UploadId
	}

	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("file does not exist: %v", filePath)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error to read [file=%v]: %v", filePath, err.Error())
	}
	chunkIndex := 1
	r := bufio.NewReader(f)
	bufferSize := 100 * 1024
	buf := make([]byte, 0, bufferSize) // 100k buffer
	fileName := path.Base(filePath)
	for {
		n, err := r.Read(buf[:cap(buf)])
		if err != nil {
			return nil, err
		}

		buf = buf[:n]
		hasMore := true
		if n < bufferSize {
			hasMore = false
		}

		// process buf
		uploadRequest(gin, accessToken, folderID, uploadId, hasMore, fileName, chunkIndex, buf)

		chunkIndex++
		if !hasMore {
			break
		}
	}
	return &dto.FolderFileChunkUploadFileResponse{}, nil
}

// Calling API Endpoint
func uploadRequest(
	gin *gin.Engine, token string, folderID int64, uploadId int64,
	hasmore bool, filename string, chunkindex int, binaryContent []byte,
) (*dto.FolderFileChunkUploadFileResponse, error) {
	var resModel dto.FolderFileChunkUploadFileResponse
	var headers = generateDefaultHeadersForAPIRequests(token)
	var hasMoreString = "1"
	if !hasmore {
		hasMoreString = "0"
	}
	params := api.FolderFileChunkUploadPathParams{
		HasMore:    hasMoreString,
		Filename:   filename,
		ChunkIndex: strconv.Itoa(chunkindex),
		UploadId:   fmt.Sprint(uploadId),
	}
	paramsMap := map[string]string{}
	bytes, _ := json.Marshal(params)
	json.Unmarshal(bytes, &paramsMap)
	for k, v := range paramsMap {
		(*headers)[k] = v
	}

	delete(*headers, "Content-Type")
	// (*headers)["Content-Type"] = "application/octet-stream"
	err := postRequestBinaryPayload(
		gin, "/api/v2/folder-files/chunk-upload-folder-file/"+strconv.Itoa(int(folderID)),
		headers, &binaryContent, &resModel,
	)
	return &resModel, err
}
