package dto

import "github.com/innovay-software/famapp-main/app/models"

type UpdateUserProfileRequest struct {
	ApiRequestBase
	Name           string `json:"name" binding:"required" validate:"required" `
	Mobile         string `json:"mobile" binding:"required" validate:"required" `
	Password       string `json:"password" binding:"omitempty" validate:"omitempty,alphanumunicode,min=6,max=30"`
	LockerPasscode string `json:"lockerPasscode" binding:"omitempty" validate:"omitempty,number,len=6" `
	AvatarUrl      string `json:"avatarUrl" binding:"omitempty" validate:"omitempty"`
}

type UpdateUserProfileResponse struct {
	ApiResponseBase `json:",squash"`
	User            *models.User `json:"userData"`
}
