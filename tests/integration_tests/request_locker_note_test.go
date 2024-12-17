package integrationTests

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/api"
	"github.com/innovay-software/famapp-main/app/dto"
)

// Calling API Endpoint
func listLockerNotes(
	gin *gin.Engine, token string,
) (
	*dto.ListLockerNotesResponse, error,
) {
	var resModel dto.ListLockerNotesResponse
	err := postRequest(
		gin, "/api/v2/locker-notes/list-notes", generateDefaultHeadersForAPIRequests(token), nil,
		&resModel,
	)
	return &resModel, err
}

// Calling API Endpoint
func createLockerNotes(
	r *gin.Engine, token, title, content string, inviteeIDs []int64,
) (
	*dto.SaveLockerNoteResponse, error,
) {
	return saveLockerNotes(r, token, 0, title, content, inviteeIDs)
}

// Calling API Endpoint
func saveLockerNotes(
	r *gin.Engine, token string, id int64, title, content string, inviteeIDs []int64,
) (
	*dto.SaveLockerNoteResponse, error,
) {
	var resModel dto.SaveLockerNoteResponse
	err := postRequest(
		r, "/api/v2/locker-notes/save-note/"+strconv.Itoa(int(id)),
		generateDefaultHeadersForAPIRequests(token),
		&api.LockerNoteSavePathJSONRequestBody{
			Title:      title,
			Content:    content,
			InviteeIds: &inviteeIDs,
		},
		&resModel,
	)
	return &resModel, err
}

// Calling API Endpoint
func deleteLockerNotes(
	r *gin.Engine, token string, id int64,
) (
	*dto.DeleteLockerNoteResponse, error,
) {
	var resModel dto.DeleteLockerNoteResponse
	err := postRequest(
		r, "/api/v2/locker-notes/delete-note/"+strconv.Itoa(int(id)),
		generateDefaultHeadersForAPIRequests(token),
		nil,
		&resModel,
	)
	return &resModel, err
}
