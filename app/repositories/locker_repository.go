package repositories

import (
	"sort"

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
	userId int64,
) (
	*[]models.LockerNote, error,
) {
	db := rp.readDBCon

	if userId <= 0 {
		return nil, errors.ApiErrorParamInvalid
	}
	var notes1 []models.LockerNote
	var notes2 []models.LockerNote

	if err := db.Order("updated_at desc").
		Preload("Owner").
		Preload("Invitees").
		Where("owner_id = ?", userId).
		Find(&notes1).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").
		Preload("Owner").
		Preload("Invitees").
		Where("id IN (SELECT note_id FROM locker_note_invitees WHERE invitee_id = ?)", userId).
		Find(&notes2).Error; err != nil {
		return nil, err
	}

	notes1 = append(notes1, notes2...)
	sort.Slice(notes1, func(i, j int) bool {
		return notes1[i].UpdatedAt.Compare(notes1[j].UpdatedAt) == 1
	})

	return &notes1, nil
}

// func (rp *lockerNoteRepo) SaveNote(
// 	user *models.User, noteId int64, title string, content string, newInviteeIds []int64,
// ) (
// 	*models.LockerNote, error,
// ) {
// 	var note models.LockerNote
// 	if noteId <= 0 {
// 		// new noteId
// 		note = models.LockerNote{OwnerID: user.ID, Title: title, Content: content}
// 		if err := rp.db.Create(&note).Error; err != nil {
// 			return nil, err
// 		}
// 		noteId = note.ID
// 	} else {
// 		// Find note record
// 		if rp.db.Find(&note, noteId).RowsAffected == 0 {
// 			note = models.LockerNote{OwnerID: user.ID, Title: title, Content: content}
// 		}

// 		// Only owner has permissions to update
// 		if note.OwnerID != user.ID {
// 			return nil, errors.ApiErrorPermissionDenied
// 		}

// 		// noteId should be found, save a version record
// 		noteVersion := models.LockerNoteVersion{
// 			Title:   note.Title,
// 			Content: note.Content,
// 			NoteID:  note.ID,
// 		}
// 		if err := SaveDbModel(&noteVersion); err != nil {
// 			return nil, err
// 		}

// 		note.Title = title
// 		note.Content = content
// 		if err := SaveDbModel(&note); err != nil {
// 			return nil, err
// 		}
// 	}
// }

func (rp *lockerNoteRepo) SyncInviteeIds(
	lockerNote *models.LockerNote, newInviteeIds *[]int64,
) error {
	db := rp.mainDBCon

	// Update invitees
	currentInviteeIds := []int64{}
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
		invitees := make([]models.LockerNoteInvitee, 0)
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
	userId int64, noteId int64,
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
		OwnerID: userId}
	return db.Delete(&note).Error
}
