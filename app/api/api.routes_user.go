/**
ApiServerInterfaceImpl - User Part
Handles all the routes related to User
*/

package api

import (
	"github.com/gin-gonic/gin"
	userHandlers "github.com/innovay-software/famapp-main/app/handlers/user"
)

// Update Profile
func (s *ApiServerInterfaceImpl) UserUpdateProfilePath(c *gin.Context, params UserUpdateProfilePathParams) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	var req UserUpdateProfilePathJSONBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := userHandlers.UpdateUserProfile(
		user, req.Name, req.Mobile, req.Password, req.LockerPasscode, req.AvatarUrl,
	)
	handleApiResponse(c, res, err)
}
