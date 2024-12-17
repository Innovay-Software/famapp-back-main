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
	c *gin.Context, user *models.User, noteId int64, title, content string, inviteeIds *[]int64,
) (
	dto.ApiResponse, error,
) {
	var lockerNote models.LockerNote
	if err := repositories.QueryDbModelByPrimaryId(
		&lockerNote, noteId,
	); err != nil {
		lockerNote = models.LockerNote{OwnerID: user.ID}
	}
	if lockerNote.OwnerID != user.ID {
		return nil, errors.ApiErrorPermissionDenied
	}

	lockerNote.Title = title
	lockerNote.Content = content
	lockerNote.ID = noteId

	if err := repositories.SaveDbModel(&lockerNote); err != nil {
		return nil, err
	}

	if err := repositories.LockerNoteRepoIns.SyncInviteeIds(
		&lockerNote, inviteeIds,
	); err != nil {
		return nil, err
	}

	res := dto.SaveLockerNoteResponse{Note: &lockerNote}
	return &res, nil
}

func DeleteLockerNoteHandler(
	c *gin.Context, user *models.User, noteId int64,
) (
	dto.ApiResponse, error,
) {
	if err := repositories.LockerNoteRepoIns.DeleteNote(user.ID, noteId); err != nil {
		return nil, errors.ApiErrorParamInvalid
	}
	return &dto.DeleteLockerNoteResponse{}, nil
}
