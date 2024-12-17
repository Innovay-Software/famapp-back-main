package dto

type DisplayFolderFileRequestUri struct {
	ApiRequestUriBase
	FolderFileId int `uri:"folderFileId" binding:"number"`
}
