package dto

import "github.com/innovay-software/famapp-main/app/models"

type ListLockerNotesResponse struct {
	ApiResponseBase `json:",squash"`
	Notes           *[]models.LockerNote `json:"notes"`
}

// type SaveLockerNoteRequestUri struct {
// 	ApiRequestUriBase
// 	NoteID int64 `uri:"noteId" binding:"number"`
// }

// type SaveLockerNoteRequest struct {
// 	ApiRequestBase
// 	Title      string  `json:"title" binding:"required"`
// 	Content    string  `json:"content" binding:"omitempty"`
// 	InviteeIDs []int64 `json:"inviteeIds" binding:"omitempty"`
// }

type SaveLockerNoteResponse struct {
	ApiResponseBase `json:",squash"`
	Note            *models.LockerNote `json:"note"`
}

// type DeleteLockerNoteRequestUri struct {
// 	ApiRequestUriBase
// 	NoteID int `uri:"noteId" binding:"number"`
// }

type DeleteLockerNoteResponse struct {
	ApiResponseBase `json:",squash"`
}
