package integrationTests

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
)

// Calling API Endpoint
func saveFolder(
	gin *gin.Engine, token string, folderID, owerID, parentID int64,
	title, cover, folderType string,
	isDefault, isPrivate bool, metadata *map[string]any, inviteeIDs []int64,
) (*dto.SaveFolderResponse, error) {
	var resModel dto.SaveFolderResponse
	if metadata == nil {
		metadata = &map[string]any{}
	}
	err := postRequest(
		gin, "/api/v2/folder/save-folder/"+strconv.Itoa(int(folderID)),
		generateDefaultHeadersForAPIRequests(token),
		&dto.SaveFolderRequest{
			OwnerID:    owerID,
			ParentID:   parentID,
			Title:      title,
			Cover:      cover,
			Type:       folderType,
			IsDefault:  isDefault,
			IsPrivate:  isPrivate,
			Metadata:   metadata,
			InviteeIds: inviteeIDs,
		}, &resModel,
	)
	return &resModel, err
}

// Calling API Endpoint
func deleteFolder(
	r *gin.Engine, token string, id int64,
) (*dto.DeleteFolderResponse, error) {
	var resModel dto.DeleteFolderResponse
	err := postRequest(
		r, "/api/v2/folder/delete-folder/"+strconv.Itoa(int(id)),
		generateDefaultHeadersForAPIRequests(token),
		nil,
		&resModel,
	)
	return &resModel, err
}
