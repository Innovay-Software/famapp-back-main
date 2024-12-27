package dto

import "github.com/innovay-software/famapp-main/app/models"

type Base64ChunkUploadFileRequest struct {
	ApiRequestBase
	Base64EncodedFileContent string `json:"base64EncodedContent" binding:"required"`
	Filename                 string `json:"filename" binding:"required"`
	HasMore                  *bool  `json:"hasMore" binding:"required"`
	ChunkedFilename          string `json:"chunkedFilename" binding:"omitempty"`
}

type Base64ChunkUploadFileResponse struct {
	ApiResponseBase `json:",squash"`
	RemoteFileId    uint64         `json:"remoteFileId"`
	Uploaded        bool           `json:"uploaded"`
	ChunkedFileName string         `json:"ChunkedFileName"`
	Document        *models.Upload `json:"document"`
	HasMore         bool           `json:"hasMore"`
}
