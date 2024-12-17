package integrationTests

import (
	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/api"
	"github.com/innovay-software/famapp-main/app/dto"
)

// Calling API Endpoint
func mobileCredentialsLogin(
	r *gin.Engine, mobile, password string,
) (
	*dto.LoginResponse, error,
) {
	var resModel dto.LoginResponse
	err := postRequest(
		r, "/api/v2/auth/mobile-login", nil, &api.AuthMobileLoginPathJSONRequestBody{
			Mobile:      mobile,
			Password:    password,
			DeviceToken: "TestingDevice",
		},
		&resModel,
	)

	return &resModel, err
}

// Calling API Endpoint
func mobileCredentialsLoginWithHeaders(
	gin *gin.Engine, mobile, password string, headers *map[string]string,
) (
	*dto.LoginResponse, error,
) {
	var resModel dto.LoginResponse
	err := postRequest(
		gin, "/api/v2/auth/mobile-login", headers, &api.AuthMobileLoginPathJSONRequestBody{
			Mobile:      mobile,
			Password:    password,
			DeviceToken: "TestingDevice",
		},
		&resModel,
	)

	return &resModel, err
}

func mobileAccessTokenLogin(
	r *gin.Engine, token string,
) (
	*dto.LoginResponse, error,
) {
	var resModel dto.LoginResponse
	err := postRequest(
		r, "/api/v2/auth/access-token-login",
		generateDefaultHeadersForAPIRequests(token), &api.AuthAccessTokenLoginPathJSONBody{
			DeviceToken: "TestingDevice",
		}, &resModel,
	)

	return &resModel, err
}
