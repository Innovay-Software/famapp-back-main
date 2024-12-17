package dto

import "github.com/innovay-software/famapp-main/app/models"

// type AdminGetMemberListRequestUri struct {
// 	ApiRequestUriBase
// 	AfterId int `uri:"afterId" binding:"number"`
// }

type AdminGetMemberListResponse struct {
	ApiResponseBase `json:",squash"`
	Users           []*models.UserMember `json:"users"`
	HasMore         bool                 `json:"hasMore"`
}

// type AdminSaveMemberRequestUri struct {
// 	ApiRequestUriBase
// 	UserId int `uri:"userId" binding:"number"`
// }

// type AdminSaveMemberRequest struct {
// 	ApiRequestBase
// 	Name           string `json:"name" binding:"omitempty" `
// 	Mobile         string `json:"mobile" binding:"omitempty" `
// 	Password       string `json:"password" binding:"omitempty" validate:"omitempty,alphanum,min=6,max=30"`
// 	LockerPasscode string `json:"lockerPasscode" binding:"omitempty" validate:"omitempty,number,len=6" `
// 	Role           string `json:"role" binding:"omitempty" `
// }

type AdminSaveMemberResponse struct {
	ApiResponseBase `json:",squash"`
	Member          *models.User `json:"member"`
}

// type AdminDeleteMemberRequestUri struct {
// 	ApiRequestUriBase
// 	UserId int `uri:"userId" binding:"number"`
// }

type AdminDeleteMemberResponse struct {
	ApiResponseBase `json:",squash"`
}
