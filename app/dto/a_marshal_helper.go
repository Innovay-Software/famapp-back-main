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
		OwnerID              uint64               `json:"ownerId"`
		ParentID             uint64               `json:"parentId"`
		Title                string               `json:"title"`
		Cover                string               `json:"cover"`
		Type                 string               `json:"type"`
		Metadata             map[string]any       `json:"metadata"`
		IsDefault            bool                 `json:"isDefault"`
		IsPrivate            bool                 `json:"isPrivate"`
		TotalFiles           uint64               `json:"totalFiles"`
		EarliestTakenOn      *time.Time           `json:"earliestTakenOn"`
		LatestTakenOn        *time.Time           `json:"latestTakenOn"`
		Invitees             []*models.UserMember `json:"invitees"`
		PopulatedSubFolders  *[]models.Folder     `json:"subFolders"`
		PopulatedLatestPosts *[]models.FolderFile `json:"latestPosts"`
	}

	// Get cover url if it doesn't start with "http"
	folder.PopulateMissingData()
	folder.Cover = utils.GetUrlPath("album-cover", folder.Cover)
	return json.Marshal(FolderMarshal(*folder))
}
