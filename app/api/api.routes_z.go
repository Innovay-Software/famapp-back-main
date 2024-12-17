package api

import (
	"reflect"

	"github.com/gin-gonic/gin"
	apiErrors "github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/innovay-software/famapp-main/config"
)

type ApiServerInterfaceImpl struct {
	// used to implement ServerInterface (generated with oapi codegen)
}

// Binds context data to request struct and validates it
func (s *ApiServerInterfaceImpl) validateInputRequest(
	c *gin.Context, req any,
) error {
	// Req must be a pointer type
	if reflect.ValueOf(req).Kind() != reflect.Ptr {
		utils.LogError("req must be a pointer")
		return apiErrors.ApiErrorSystem
	}
	// Binds to req struct
	if err := c.Bind(req); err != nil {
		return err
	}
	// Validates req based on its tags
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

// Get authenticated user
func getAuthenticatedUser(
	c *gin.Context,
) (
	*models.User, error,
) {
	// AuthenticatedUser should be injected by middleware
	userObject, exists := c.Get(config.AuthenticatedUserHeaderKey)
	if !exists {
		return nil, apiErrors.ApiErrorToken
	}
	user, ok := userObject.(*models.User)
	if !ok {
		return nil, apiErrors.ApiErrorToken
	}

	return user, nil
}

// Get authenticated admin user from middleware
func getAuthenticatedAdmin(c *gin.Context) (*models.User, error) {
	user, err := getAuthenticatedUser(c)
	if err != nil {
		return nil, err
	}
	if !user.IsAdmin() {
		return nil, apiErrors.ApiErrorRequiresAdmin
	}
	return user, nil
}
