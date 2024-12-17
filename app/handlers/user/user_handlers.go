package user

import (
	"github.com/innovay-software/famapp-main/app/dto"
	apiErrors "github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUserProfile(
	user *models.User, name, mobile, password, lockerPasscode, avatarUrl *string,
) (
	dto.ApiResponse, error,
) {
	if mobile != nil && *mobile != user.Mobile {
		// Change mobile number
		_, err := repositories.UserRepoIns.FindUserByField("mobile", *mobile)
		if err == nil {
			// Is able to find a user with that new mobile
			return nil, apiErrors.ApiErrorDuplicateMobile
		}
		user.Mobile = *mobile
	}
	if name != nil {
		user.Name = *name
	}
	if password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, apiErrors.ApiErrorSystem
		}
		user.Password = string(hashedPassword)
	}
	if lockerPasscode != nil {
		user.LockerPasscode = *lockerPasscode
	}
	if avatarUrl != nil {
		user.SetAvatarUrl(*avatarUrl)
	}

	if err := repositories.UserRepoIns.SaveUser(user); err != nil {
		return nil, apiErrors.ApiErrorSystem
	}
	return &dto.UpdateUserProfileResponse{
		User: user,
	}, nil
}
