package dto

import "github.com/innovay-software/famapp-main/app/models"

type SaveFolderRequestUri struct {
	ApiRequestUriBase
	FolderId int64 `uri:"folderId" binding:"number"`
}

type SaveFolderRequest struct {
	ApiRequestBase
	OwnerID    int64           `json:"ownerId" binding:"required"`
	ParentID   int64           `json:"parentId" binding:"omitempty"`
	Title      string          `json:"title" binding:"required"`
	Cover      string          `json:"cover" binding:"omitempty"`
	Type       string          `json:"type" binding:"required"`
	IsDefault  bool            `json:"isDefault" binding:"omitempty"`
	IsPrivate  bool            `json:"isPrivate" binding:"omitempty"`
	Metadata   *map[string]any `json:"metadata" binding:"omitempty"`
	InviteeIds []int64         `json:"inviteeIds" binding:"omitempty"`
}

type SaveFolderResponse struct {
	ApiResponseBase `json:",squash"`
	Folder          *models.Folder `json:"folder"`
}

type DeleteFolderRequestUri struct {
	ApiRequestUriBase
	FolderId int64 `uri:"folderId" binding:"required,number"`
}

type DeleteFolderResponse struct {
	ApiResponseBase `json:",squash"`
}
