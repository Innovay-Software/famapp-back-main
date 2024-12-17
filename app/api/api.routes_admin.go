/**
ApiServerInterfaceImpl - Admin Part
Handles all the admin related routes
*/

package api

import (
	"github.com/gin-gonic/gin"
	adminHandlers "github.com/innovay-software/famapp-main/app/handlers/admin"
)

// Admin adds an user (create, update)
func (s *ApiServerInterfaceImpl) AdminAddUserPath(
	c *gin.Context, params AdminAddUserPathParams,
) {
	// Authenticate caller
	admin, err := getAuthenticatedAdmin(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	var req AdminAddUserPathJSONRequestBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := adminHandlers.AdminAddMember(
		c, admin, req.Name, req.Mobile,
		string(req.Role), req.Password, req.LockerPasscode, req.FamilyId,
	)
	handleApiResponse(c, res, err)
}

// Admin list users
func (s *ApiServerInterfaceImpl) AdminListUsersPath(
	c *gin.Context, afterId int64, params AdminListUsersPathParams,
) {
	// Authenticate caller
	admin, err := getAuthenticatedAdmin(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request body

	// Handle request
	res, err := adminHandlers.AdminGetMemberListHandler(c, admin, afterId)
	handleApiResponse(c, res, err)
}

// Admin delete a user
func (s *ApiServerInterfaceImpl) AdminDeleteUserPath(
	c *gin.Context, uuid string, params AdminDeleteUserPathParams,
) {
	// Authenticate caller
	admin, err := getAuthenticatedAdmin(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request body

	// Handle request
	res, err := adminHandlers.AdminDeleteMember(c, admin, uuid)
	handleApiResponse(c, res, err)
}

// Admin saves an user (create, update)
func (s *ApiServerInterfaceImpl) AdminSaveUserPath(
	c *gin.Context, userId int64, params AdminSaveUserPathParams,
) {
	// Authenticate caller
	admin, err := getAuthenticatedAdmin(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	var req AdminSaveUserPathJSONRequestBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := adminHandlers.AdminUpdateMember(
		c, admin, userId, req.Name, req.Mobile,
		string(req.Role), req.Password, req.LockerPasscode, req.FamilyId,
	)
	handleApiResponse(c, res, err)
}
