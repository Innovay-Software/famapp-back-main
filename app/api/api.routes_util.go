/**
ApiServerInterfaceImpl - Utils Part
Handles all the routes related to Utils
*/

package api

import (
	"github.com/gin-gonic/gin"
	utilHandlers "github.com/innovay-software/famapp-main/app/handlers/util"
)

// Get config based on key provided
func (s *ApiServerInterfaceImpl) UtilConfigPath(
	c *gin.Context, configKey string, params UtilConfigPathParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request body

	// Handle request
	res, err := utilHandlers.GetConfig(c, user, configKey)
	handleApiResponse(c, res, err)
}

// Display user avatar
func (s *ApiServerInterfaceImpl) UtilDisplayUserAvatarPath(c *gin.Context, userId int64, params UtilDisplayUserAvatarPathParams) {
	// Authenticate caller
	// pass - public api

	// Validate Request
	// pass - no request body

	// Handle request
	err := utilHandlers.UserAvatar(c, userId)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}
}

// Check for front end updates
func (s *ApiServerInterfaceImpl) UtilCheckForMobileUpdatePath(
	c *gin.Context, clientOs string, clientVersion string, params UtilCheckForMobileUpdatePathParams,
) {
	// Authenticate caller
	// pass - public api

	// Validate Request
	// pass - no request body

	// Handle request
	res, err := utilHandlers.CheckForMobileUpdate(c, clientOs, clientVersion)
	handleApiResponse(c, res, err)
}

// Ping server for health
func (s *ApiServerInterfaceImpl) UtilPingPath(
	c *gin.Context, params UtilPingPathParams,
) {
	// Authenticate caller
	// pass - public api

	// Validate Request
	// pass - no request body

	// Handle request
	res, err := utilHandlers.Ping(c)
	handleApiResponse(c, res, err)
}
