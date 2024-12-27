package models

type FolderInvitee struct {
	BaseDbModel
	FolderID  uint64 `gorm:"column:folder_id" json:"folderId"`
	InviteeID uint64 `gorm:"column:invitee_id" json:"inviteeId"`
}

func (FolderInvitee) TableName() string {
	return "folder_invitees"
}
