package integrationTests

import (
	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
)

// Create a new user
func createUser(
	gin *gin.Engine, name, mobile string, role string,
) (
	*dto.AdminSaveMemberResponse, error,
) {
	// Admin login
	res1, err1 := mobileCredentialsLogin(gin, "1234567890", "123456")
	if err1 != nil {
		return nil, err1
	}
	adminToken := res1.AccessToken
	password := "123456"
	passcode := "123456"

	return adminAddMemberAPI(gin, adminToken, name, mobile, role, password, passcode)
}
