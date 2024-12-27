package repositories

import (
	"slices"

	"github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/utils"

	"gorm.io/gorm"
)

type lockerNoteRepo struct {
	mainDBCon *gorm.DB
	readDBCon *gorm.DB
	rd        *redisRepo
}

func (rp *lockerNoteRepo) FindAllNotes(
	userId uint64,
) (
	*[]models.LockerNote, error,
) {
	db := rp.readDBCon

	if userId <= 0 {
		return nil, errors.ApiErrorParamInvalid
	}
	var notes []models.LockerNote

	if err := db.Order("updated_at desc").
		Preload("Owner").
		Preload("Invitees").
		Where("id IN (SELECT note_id FROM locker_note_invitees WHERE invitee_id = ?)", userId).
		Find(&notes).Error; err != nil {
		return nil, err
	}

	return &notes, nil
}

func (rp *lockerNoteRepo) SaveNote(
	user *models.User, noteId uint64, title, content string, newInviteeIds *[]uint64,
) (
	*models.LockerNote, error,
) {
	if !slices.Contains(*newInviteeIds, user.ID) {
		*newInviteeIds = append(*newInviteeIds, user.ID)
	}

	var lockerNote *models.LockerNote
	if noteId > 0 {
		currentLockerNote := &models.LockerNote{}
		if err := QueryDbModelByPrimaryId(
			currentLockerNote, noteId,
		); err != nil {
			utils.LogError("err = ", err)
			lockerNote = nil
		} else {
			lockerNote = currentLockerNote
		}
	}
	if lockerNote == nil {
		utils.LogError("create new lockernote")
		lockerNote = &models.LockerNote{OwnerID: user.ID}
	}

	if lockerNote.OwnerID != user.ID {
		return nil, errors.ApiErrorPermissionDenied
	}

	lockerNote.Title = title
	lockerNote.Content = content
	if noteId > 0 {
		lockerNote.ID = noteId
	}

	if err := SaveDbModel(lockerNote); err != nil {
		return nil, err
	}

	if err := LockerNoteRepoIns.SyncInviteeIds(
		lockerNote, newInviteeIds,
	); err != nil {
		return nil, err
	}

	return lockerNote, nil
}

func (rp *lockerNoteRepo) SyncInviteeIds(
	lockerNote *models.LockerNote, newInviteeIds *[]uint64,
) error {
	db := rp.mainDBCon

	// Update invitees
	currentInviteeIds := []uint64{}
	db.Model(&models.LockerNoteInvitee{}).
		Where("note_id = ?", lockerNote.ID).
		Pluck("invitee_id", &currentInviteeIds)

	deleteInviteeIds := utils.SliceLeftExcludeRight(&currentInviteeIds, newInviteeIds)
	insertInviteeIds := utils.SliceLeftExcludeRight(newInviteeIds, &currentInviteeIds)

	if len(*deleteInviteeIds) > 0 {
		if err := db.Where("note_id = ?", lockerNote.ID).
			Where("invitee_id in ?", *deleteInviteeIds).
			Delete(&models.LockerNoteInvitee{}).Error; err != nil {
			return err
		}
	}
	if len(*insertInviteeIds) > 0 {
		invitees := []models.LockerNoteInvitee{}
		for _, item := range *insertInviteeIds {
			invitees = append(invitees, models.LockerNoteInvitee{NoteID: lockerNote.ID, InviteeID: item})
		}
		if err := db.Create(invitees).Error; err != nil {
			return err
		}
	}

	return nil
}

func (rp *lockerNoteRepo) DeleteNote(
	userId, noteId uint64,
) error {
	db := rp.mainDBCon

	if userId <= 0 || noteId <= 0 {
		return errors.ApiErrorParamInvalid
	}
	note := models.LockerNote{
		BaseModelSoftDelete: models.BaseModelSoftDelete{
			BaseDbModel: models.BaseDbModel{
				ID: noteId,
			},
		},
		OwnerID: userId,
	}
	return db.Delete(&note).Error
}
