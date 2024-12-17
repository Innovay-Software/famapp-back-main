package models

type LockerNote struct {
	BaseModelSoftDelete
	Title    string        `gorm:"column:title" json:"title"`
	Content  string        `gorm:"column:content" json:"content"`
	OwnerID  int64         `gorm:"column:owner_id" json:"ownerId"`
	Owner    *User         `gorm:"foreignKey:OwnerID" json:"owner"`
	Invitees []*UserMember `gorm:"many2many:locker_note_invitees;foreignKey:ID;joinForeignKey:NoteID;References:ID;joinReferences:InviteeID" json:"invitees"`
}

func (LockerNote) TableName() string {
	return "locker_notes"
}

type LockerNoteVersion struct {
	BaseDbModel
	Title   string `gorm:"column:title" json:"title"`
	Content string `gorm:"column:content" json:"content"`
	NoteID  int64  `gorm:"column:note_id" json:"noteId"`
}

func (LockerNoteVersion) TableName() string {
	return "locker_note_versions"
}

type LockerNoteInvitee struct {
	BaseDbModel
	NoteID    int64 `gorm:"column:note_id" json:"noteId"`
	InviteeID int64 `gorm:"column:invitee_id" json:"inviteeId"`
}

func (LockerNoteInvitee) TableName() string {
	return "locker_note_invitees"
}
