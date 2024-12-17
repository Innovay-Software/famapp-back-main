package integrationTests

import (
	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/api"
	"github.com/innovay-software/famapp-main/app/dto"
)

// Calling API Endpoint
func updateUserProfile(
	r *gin.Engine, token string, name, mobile, password, lockerPasscode, avatarUrl *string,
) (
	*dto.UpdateUserProfileResponse, error,
) {
	var resModel dto.UpdateUserProfileResponse
	err := postRequest(
		r, "/api/v2/user/update-profile",
		generateDefaultHeadersForAPIRequests(token),
		&api.UserUpdateProfilePathJSONBody{
			Name:           name,
			Mobile:         mobile,
			Password:       password,
			LockerPasscode: lockerPasscode,
			AvatarUrl:      avatarUrl,
		},
		&resModel,
	)
	return &resModel, err
}
