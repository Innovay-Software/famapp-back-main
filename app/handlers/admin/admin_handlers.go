package admin

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/gin-gonic/gin"

	"github.com/innovay-software/famapp-main/app/dto"
	apiErrors "github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
)

// Add member (create or update)
func AdminAddMember(
	c *gin.Context, admin *models.User,
	name, mobile, role, password, lockerPasscode string, familyId uint64,
) (
	dto.ApiResponse, error,
) {
	h := md5.New()
	h.Write([]byte(lockerPasscode))
	newUser := models.User{
		FamilyID:       familyId,
		Name:           name,
		Mobile:         mobile,
		Role:           role,
		LockerPasscode: hex.EncodeToString(h.Sum(nil)),
	}
	newUser.SetPassword(password)
	if err := repositories.UserRepoIns.CreateUser(&newUser); err != nil {
		return nil, err
	}

	return &dto.AdminSaveMemberResponse{Member: &newUser}, nil
}

func AdminGetMemberListHandler(
	c *gin.Context, admin *models.User, afterId uint64,
) (
	dto.ApiResponse, error,
) {
	pageSize := 100
	members, _ := repositories.UserRepoIns.FindMemberList(
		pageSize, afterId,
	)

	res := dto.AdminGetMemberListResponse{
		Users:   members,
		HasMore: len(members) == pageSize,
	}
	return &res, nil
}

// Update member (Member must already exist)
func AdminUpdateMember(
	c *gin.Context, admin *models.User,
	targetUserId uint64, name, mobile, role string, password, lockerPasscode *string, familyId *uint64,
) (
	dto.ApiResponse, error,
) {
	if targetUserId <= 0 {
		return nil, apiErrors.ApiErrorParamInvalid
	}

	var targetUser models.User = models.User{}
	err := repositories.QueryDbModelByPrimaryId(&targetUser, targetUserId)
	if err != nil {
		return nil, err
	}

	// Update password
	if password != nil && *password != "" {
		targetUser.SetPassword(*password)
	}
	if lockerPasscode != nil {
		h := md5.New()
		h.Write([]byte(*lockerPasscode))
		targetUser.LockerPasscode = hex.EncodeToString(h.Sum(nil))
	}
	if name != "" {
		targetUser.Name = name
	}
	if mobile != "" {
		targetUser.Mobile = mobile
	}
	if role != "" {
		targetUser.Role = role
	}
	if familyId != nil {
		targetUser.FamilyID = uint64(*familyId)
	}

	if err := repositories.UserRepoIns.SaveUser(&targetUser); err != nil {
		return nil, err
	}

	return &dto.AdminSaveMemberResponse{Member: &targetUser}, nil
}

// Delete member
func AdminDeleteMember(
	c *gin.Context, admin *models.User, targetUUID string,
) (
	dto.ApiResponse, error,
) {
	repositories.UserRepoIns.DeleteUser("uuid", targetUUID)
	return &dto.AdminDeleteMemberResponse{}, nil
}
