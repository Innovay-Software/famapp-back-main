package dto

import (
	"encoding/json"
	"time"

	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/utils"
)

func marshalFolder(folder *models.Folder) ([]byte, error) {
	// Define a temporary struct to hold the marshalled data
	type FolderMarshal struct {
		models.BaseDbModel
		OwnerID              int64                `json:"ownerId"`
		ParentID             int64                `json:"parentId"`
		Title                string               `json:"title"`
		Cover                string               `json:"cover"`
		Type                 string               `json:"type"`
		Metadata             map[string]any       `json:"metadata"`
		IsDefault            bool                 `json:"isDefault"`
		IsPrivate            bool                 `json:"isPrivate"`
		TotalFiles           int64                `json:"totalFiles"`
		EarliestShotAt       *time.Time           `json:"earliestShotAt"`
		LatestShotAt         *time.Time           `json:"latestShotAt"`
		Invitees             []*models.UserMember `json:"invitees"`
		PopulatedSubFolders  *[]models.Folder     `json:"subFolders"`
		PopulatedLatestPosts *[]models.FolderFile `json:"latestPosts"`
	}

	// Get cover url if it doesn't start with "http"
	folder.PopulateMissingData()
	folder.Cover = utils.GetUrlPath("album-cover", folder.Cover)
	return json.Marshal(FolderMarshal(*folder))
}
