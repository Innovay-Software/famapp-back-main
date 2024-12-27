package lockerNote

import (
	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
	"github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
)

func ListLockerNotesHandler(
	c *gin.Context, user *models.User,
) (
	dto.ApiResponse, error,
) {
	notes, err := repositories.LockerNoteRepoIns.FindAllNotes(user.ID)
	if err != nil {
		return nil, errors.ApiErrorSystem
	}

	res := dto.ListLockerNotesResponse{
		Notes: notes,
	}

	return &res, nil
}

// Save a locker note
func SaveLockerNoteHandler(
	c *gin.Context, user *models.User, noteId uint64, title, content string, inviteeIds *[]uint64,
) (
	dto.ApiResponse, error,
) {

	lockerNote, err := repositories.LockerNoteRepoIns.SaveNote(
		user, noteId, title, content, inviteeIds,
	)
	if err != nil {
		return nil, err
	}

	res := dto.SaveLockerNoteResponse{Note: lockerNote}
	return &res, nil
}

func DeleteLockerNoteHandler(
	c *gin.Context, user *models.User, noteId uint64,
) (
	dto.ApiResponse, error,
) {
	if err := repositories.LockerNoteRepoIns.DeleteNote(user.ID, noteId); err != nil {
		return nil, errors.ApiErrorParamInvalid
	}
	return &dto.DeleteLockerNoteResponse{}, nil
}
