package models

type FolderInvitee struct {
	BaseDbModel
	FolderID  int64 `gorm:"column:folder_id" json:"folderId"`
	InviteeID int64 `gorm:"column:invitee_id" json:"inviteeId"`
}

func (FolderInvitee) TableName() string {
	return "folder_invitees"
}
