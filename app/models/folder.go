package models

import (
	"encoding/json"
	"time"

	"github.com/innovay-software/famapp-main/app/services"
	"github.com/innovay-software/famapp-main/app/utils"
)

type Folder struct {
	BaseDbModel
	OwnerID              int64          `gorm:"column:owner_id" json:"ownerId"`
	ParentID             int64          `gorm:"column:parent_id" json:"parentId"`
	Title                string         `gorm:"column:title" json:"title"`
	Cover                string         `gorm:"column:cover" json:"cover"`
	Type                 string         `gorm:"column:type" json:"type"`
	Metadata             map[string]any `gorm:"column:metadata;serializer:json" json:"metadata"`
	IsDefault            bool           `gorm:"column:is_default; default:0" json:"isDefault"`
	IsPrivate            bool           `gorm:"column:is_private; default:0" json:"isPrivate"`
	TotalFiles           int64          `gorm:"column:total_files; default:0" json:"totalFiles"`
	EarliestShotAt       *time.Time     `gorm:"column:earliest_shot_at; null" json:"earliestShotAt"`
	LatestShotAt         *time.Time     `gorm:"column:latest_shot_at; null" json:"latestShotAt"`
	Invitees             []*UserMember  `gorm:"many2many:folder_invitees;foreignKey:ID;joinForeignKey:FolderID;References:ID;joinReferences:InviteeID" json:"invitees"`
	PopulatedSubFolders  *[]Folder      `gorm:"-" json:"subFolders"`
	PopulatedLatestPosts *[]FolderFile  `gorm:"-" json:"latestPosts"`
}

func (Folder) TableName() string {
	return "folders"
}

func (f *Folder) PopulateMissingData() {
	const fileLimit = 100
	db := services.GetReadDBConnection()

	if f.PopulatedLatestPosts == nil {
		var files []FolderFile
		db.Limit(fileLimit).
			Where("folder_id = ?", f.ID).
			Order("shot_at desc").Find(&files)
		f.PopulatedLatestPosts = &files
	}

	if f.PopulatedSubFolders == nil {
		var folders []Folder
		db.Where("parent_id = ?", f.ID).Find(&folders)
		f.PopulatedSubFolders = &folders
	}
}

func (f *Folder) MarshalJSON() ([]byte, error) {
	// Define a temporary struct to hold the marshalled data
	type FolderMarshal struct {
		BaseDbModel
		OwnerID              int64          `json:"ownerId"`
		ParentID             int64          `json:"parentId"`
		Title                string         `json:"title"`
		Cover                string         `json:"cover"`
		Type                 string         `json:"type"`
		Metadata             map[string]any `json:"metadata"`
		IsDefault            bool           `json:"isDefault"`
		IsPrivate            bool           `json:"isPrivate"`
		TotalFiles           int64          `json:"totalFiles"`
		EarliestShotAt       *time.Time     `json:"earliestShotAt"`
		LatestShotAt         *time.Time     `json:"latestShotAt"`
		Invitees             []*UserMember  `json:"invitees"`
		PopulatedSubFolders  *[]Folder      `json:"subFolders"`
		PopulatedLatestPosts *[]FolderFile  `json:"latestPosts"`
	}

	// Get cover url if it doesn't start with "http"
	f.PopulateMissingData()
	f.Cover = utils.GetUrlPath("album-cover", f.Cover)
	return json.Marshal(FolderMarshal(*f))
}

func (f *Folder) ToClientMap() map[string]any {
	type FolderForClient struct {
		BaseDbModel
		OwnerID              int64          `json:"ownerId"`
		ParentID             int64          `json:"parentId"`
		Title                string         `json:"title"`
		Cover                string         `json:"cover"`
		Type                 string         `json:"type"`
		Metadata             map[string]any `json:"metadata"`
		IsDefault            bool           `json:"isDefault"`
		IsPrivate            bool           `json:"isPrivate"`
		TotalFiles           int64          `json:"totalFiles"`
		EarliestShotAt       *time.Time     `json:"earliestShotAt"`
		LatestShotAt         *time.Time     `json:"latestShotAt"`
		Invitees             []*UserMember  `json:"invitees"`
		PopulatedSubFolders  *[]Folder      `json:"subFolders"`
		PopulatedLatestPosts *[]FolderFile  `json:"latestPosts"`
	}

	// Get cover url if it doesn't start with "http"
	f.PopulateMissingData()
	f.Cover = utils.GetUrlPath("album-cover", f.Cover)
	t := FolderForClient(*f)

	res := map[string]any{}
	jsonBytes, _ := json.Marshal(t)
	json.Unmarshal(jsonBytes, &res)
	return res
}
