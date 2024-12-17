package dto

import (
	"github.com/innovay-software/famapp-main/app/models"
)

// type FolderFileChunkUploadFileRequestUri struct {
// 	ApiRequestUriBase
// 	FolderId int `uri:"folderId" binding:"number"`
// }

type FolderFileChunkUploadFileResponse struct {
	ApiResponseBase `json:",squash"`
}

type FolderFileGetChunkUploadFileIdResponse struct {
	ApiResponseBase `json:",squash"`
	UploadId        int64 `json:"uploadId"`
}

// type GetFolderFilesBeforeIdRequestUri struct {
// 	ApiRequestUriBase
// 	FolderId  int64  `uri:"folderId" binding:"number" validate:"required,gt=0"`
// 	PivotDate string `uri:"pivotDate" binding:"required" validate:"required"`
// }

// type GetFolderFilesBeforeIdRequest struct {
// 	ApiRequestBase
// 	BeforeShotAt string `json:"beforeShotAt" binding:"omitempty"`
// }

type GetFolderFilesBeforeShotAtResponse struct {
	ApiResponseBase `json:",squash"`
	FolderFiles     *[]models.FolderFile `json:"folderFiles"`
	Folder          *models.Folder       `json:"folder"`
	HasMore         bool                 `json:"hasMore"`
}

// type GetFolderFilesAfterIdRequestUri struct {
// 	ApiRequestUriBase
// 	FolderId  int64  `uri:"folderId" binding:"number"`
// 	PivotDate string `uri:"pivotDate" binding:"required"`
// }

// type GetFolderFilesAfterIdRequest struct {
// 	ApiRequestBase
// 	AfterShotAt string `json:"afterShotAt" binding:"required"`
// }

type GetFolderFilesAfterShotAtResponse struct {
	ApiResponseBase `json:",squash"`
	FolderFiles     *[]models.FolderFile `json:"folderFiles"`
	Folder          *models.Folder       `json:"folder"`
	HasMore         bool                 `json:"hasMore"`
}

// type UpdateSingleFolderFileRequestUri struct {
// 	ApiRequestUriBase
// 	FolderFileId int64 `uri:"folderFileId" binding:"required,number"`
// }

// type UpdateSingleFolderFileRequest struct {
// 	ApiRequestBase
// 	Remark      string `json:"remark" binding:"omitempty"`
// 	IsPrivate   bool   `json:"isPrivate" binding:"required"`
// 	NewFolderId int64  `json:"folderId" binding:"number"`
// }

type UpdateSingleFolderFileResponse struct {
	ApiResponseBase `json:",squash"`
}

// type UpdateMultipleFolderFileRequest struct {
// 	ApiRequestBase
// 	FolderFileIds     []int64    `json:"fileIds" binding:"omitempty"`
// 	NewFolderId       int64      `json:"newFolderId" binding:"omitempty"`
// 	NewShotAtDateTime *time.Time `json:"newShotAtDateTime" binding:"omitempty"`
// }

type UpdateMultipleFolderFileResponse struct {
	ApiResponseBase `json:",squash"`
}

// type DeleteFolderFileRequest struct {
// 	ApiRequestBase
// 	FolderFileIds []int64 `json:"fileIds" binding:"required"`
// }

type DeleteFolderFileResponse struct {
	ApiResponseBase `json:",squash"`
}
