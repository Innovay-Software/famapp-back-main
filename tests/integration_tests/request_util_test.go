package integrationTests

import (
	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
)

// Calling API Endpoint
func utilGetConfigAPI(
	r *gin.Engine, token, configKey string,
) (
	*dto.GetConfigResponse, error,
) {
	var resModel dto.GetConfigResponse
	err := postRequest(
		r, "/api/v2/util/config/"+configKey,
		generateDefaultHeadersForAPIRequests(token),
		nil,
		&resModel,
	)
	return &resModel, err
}

func utilPingAPI(
	r *gin.Engine,
) (
	*dto.PingResponse, error,
) {
	var resModel dto.PingResponse
	err := getRequest(
		r, "/api/v2/util/ping",
		nil,
		&resModel,
	)
	return &resModel, err
}

func utilCheckForMobileUpdate(
	r *gin.Engine, clientOS, clientVersion string,
) (
	*dto.CheckForUpdateResponse, error,
) {
	var resModel dto.CheckForUpdateResponse
	err := getRequest(
		r, "/api/v2/util/check-for-mobile-update/"+clientOS+"/"+clientVersion,
		nil,
		&resModel,
	)
	return &resModel, err
}
