/**
ApiServerInterfaceImpl - Auth Part
Handles all the login / registration routes
*/

package api

import (
	"github.com/gin-gonic/gin"
	authHandlers "github.com/innovay-software/famapp-main/app/handlers/auth"
)

// Mobile login
func (s *ApiServerInterfaceImpl) AuthMobileLoginPath(
	c *gin.Context, params AuthMobileLoginPathParams,
) {
	// Authenticate caller
	// pass - user not logged in

	// Validate Request
	var req AuthMobileLoginPathJSONBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := authHandlers.MobileLoginHandler(
		c, req.Mobile, req.Password, req.DeviceToken,
	)
	handleApiResponse(c, res, err)
}

// AccessToken login
func (s *ApiServerInterfaceImpl) AuthAccessTokenLoginPath(c *gin.Context, params AuthAccessTokenLoginPathParams) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	var req AuthAccessTokenLoginPathJSONBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := authHandlers.AccessTokenLoginHandler(
		c, user, req.DeviceToken,
	)
	handleApiResponse(c, res, err)
}
