package integrationTests

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/api"
	"github.com/innovay-software/famapp-main/app/dto"
)

func getLatestFolderFilesAfterTimestamp(
	r *gin.Engine, token string, folderId int64, pivotDate string, microTimestamp int64,
) (
	*dto.GetFolderFilesBeforeTakenOnResponse, error,
) {
	var resModel dto.GetFolderFilesBeforeTakenOnResponse
	err := postRequest(
		r, fmt.Sprintf("/api/v2/folder-files/get-folder-files-after-micro-timestamp/%v/%v/%v", folderId, pivotDate, microTimestamp), generateDefaultHeadersForAPIRequests(token), nil, &resModel,
	)
	return &resModel, err
}

func getLatestFolderFilesBeforeTimestamp(
	r *gin.Engine, token string, folderId int64, pivotDate string, microTimestamp int64,
) (
	*dto.GetFolderFilesBeforeTakenOnResponse, error,
) {
	var resModel dto.GetFolderFilesBeforeTakenOnResponse
	err := postRequest(
		r, fmt.Sprintf("/api/v2/folder-files/get-folder-files-before-micro-timestamp/%v/%v/%v", folderId, pivotDate, microTimestamp), generateDefaultHeadersForAPIRequests(token), nil, &resModel,
	)
	return &resModel, err
}

func updateSingleFolderFile(
	r *gin.Engine, token string, folderFileId int64, remark *string, isPrivate *bool,
) (
	*dto.UpdateSingleFolderFileResponse, error,
) {
	var resModel dto.UpdateSingleFolderFileResponse

	err := postRequest(
		r, "/api/v2/folder-files/update-single-folder-file", generateDefaultHeadersForAPIRequests(token),
		&api.FolderFileUpdateSingleFolderFilePathJSONRequestBody{
			FolderFileId: folderFileId,
			IsPrivate:    isPrivate,
			Remark:       remark,
		}, &resModel,
	)
	return &resModel, err
}
