package integrationTests

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/api"
	"github.com/innovay-software/famapp-main/app/dto"
)

// Calling API Endpoint
func adminListMembersAPI(
	r *gin.Engine, token string,
) (
	*dto.AdminGetMemberListResponse, error,
) {
	var resModel dto.AdminGetMemberListResponse
	err := postRequest(
		r, "/api/v2/admin/list-all-users/0", generateDefaultHeadersForAPIRequests(token), nil,
		&resModel,
	)
	return &resModel, err
}

// Calling API Endpoint
func adminAddMemberAPI(
	r *gin.Engine, token, name, mobile, role, password, lockerPasscode string,
) (
	*dto.AdminSaveMemberResponse, error,
) {
	var resModel dto.AdminSaveMemberResponse
	err := postRequest(
		r, "/api/v2/admin/add-user", generateDefaultHeadersForAPIRequests(token),
		&api.AdminAddUserPathJSONRequestBody{
			Mobile:         mobile,
			Name:           name,
			Password:       password,
			LockerPasscode: lockerPasscode,
			Role:           role,
			FamilyId:       1,
		},
		&resModel,
	)
	return &resModel, err
}

// Calling API Endpoint
func adminUpdateMemberAPI(
	r *gin.Engine, token string, userId uint64, name, mobile, role string,
	password, lockerPasscode *string,
) (
	*dto.AdminSaveMemberResponse, error,
) {
	var resModel dto.AdminSaveMemberResponse
	familyId := 1
	err := postRequest(
		r, "/api/v2/admin/update-user/"+strconv.Itoa(int(userId)), generateDefaultHeadersForAPIRequests(token),
		&api.AdminSaveUserPathJSONRequestBody{
			LockerPasscode: lockerPasscode,
			Mobile:         mobile,
			Name:           name,
			Password:       password,
			Role:           role,
			FamilyId:       &familyId,
		},
		&resModel,
	)
	return &resModel, err
}

// Calling API Endpoint
func adminDeleteMemberAPI(
	r *gin.Engine, token string, userUuid string,
) (
	*dto.AdminDeleteMemberResponse, error,
) {
	var resModel dto.AdminDeleteMemberResponse
	err := postRequest(
		r, "/api/v2/admin/delete-user/"+userUuid, generateDefaultHeadersForAPIRequests(token),
		nil, &resModel,
	)
	return &resModel, err
}
