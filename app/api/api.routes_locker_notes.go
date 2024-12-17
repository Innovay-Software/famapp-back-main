/**
ApiServerInterfaceImpl - Locker Notes Part
Handles all the routes related to locker notes
*/

package api

import (
	"github.com/gin-gonic/gin"
	lockerNoteHandlers "github.com/innovay-software/famapp-main/app/handlers/lockerNote"
)

// List LockerNotes
func (s *ApiServerInterfaceImpl) LockerNoteListPath(c *gin.Context, params LockerNoteListPathParams) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request bdy

	// Handle request
	res, err := lockerNoteHandlers.ListLockerNotesHandler(
		c, user,
	)
	handleApiResponse(c, res, err)
}

// Delete LockerNote
func (s *ApiServerInterfaceImpl) LockerNoteDeletePath(c *gin.Context, noteId int64, params LockerNoteDeletePathParams) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request bdy

	// Handle request
	res, err := lockerNoteHandlers.DeleteLockerNoteHandler(
		c, user, noteId,
	)
	handleApiResponse(c, res, err)
}

// Save LockerNote
func (s *ApiServerInterfaceImpl) LockerNoteSavePath(c *gin.Context, noteId int64, params LockerNoteSavePathParams) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	var req LockerNoteSavePathJSONRequestBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := lockerNoteHandlers.SaveLockerNoteHandler(
		c, user, noteId, req.Title, req.Content, req.InviteeIds,
	)
	handleApiResponse(c, res, err)
}
