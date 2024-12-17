package integrationTests

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app"
	"github.com/innovay-software/famapp-main/tests"
)

// Workflow1 tests for user authentication, user profile
func TestWorkflow1(t *testing.T) {
	if !runWorkflow1 {
		return
	}

	_, b, _, _ := runtime.Caller(0)
	projDir := filepath.Dir(b)
	r, _ := app.InitApiIntegrationTestServer(fmt.Sprintf("%s/../..", projDir))

	workflow1A_Test(t, r)
	workflow1B_Test(t, r)
	workflow1D_Test(t, r)

	// Reset database
	resetDatabase()
}

// Workflow1A - Login and Refresh tokens
func workflow1A_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow1A_Test: Login and Refresh tokens")

	// Reset database
	resetDatabase()

	// UserID1 login
	{
		loginRes1, err := mobileCredentialsLogin(r, superAdminMobile, superAdminPassword)
		tests.AssertNil(t, err)
		tests.AssertEqual(t, loginRes1.User.ID, superAdminID)
		tests.AssertEqual(t, loginRes1.User.Mobile, superAdminMobile)
		user1AccessToken1 := loginRes1.AccessToken
		tests.AssertNotEqual(t, user1AccessToken1, "")
	}

	logWorkflowSuccess("Workflow1A_Test: ended")
}

// Workflow1B - User profile update
func workflow1B_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow1B_Test: User profile update")

	// Reset database
	resetDatabase()

	// User login
	loginRes1, _ := mobileCredentialsLogin(r, superAdminMobile, superAdminPassword)
	user1AccessToken1 := loginRes1.AccessToken

	// Upload file
	uploadedFilePath, fileUrl, err := utilBase64ChunkUploadFileWrapper(t, r, user1AccessToken1, "../files/sample-image.jpeg")
	tests.AssertNil(t, err)
	tests.AssertNotEqual(t, uploadedFilePath, "")
	defer os.Remove(uploadedFilePath)

	// Update user profile
	res1, err1 := updateUserProfile(r, user1AccessToken1, nil, nil, nil, nil, &fileUrl)
	tests.AssertNil(t, err1)
	tests.AssertEqual(t, fileUrl, res1.User.Avatar)

	logWorkflowSuccess("Workflow1B_Test: ended")
}

// Workflow1D: Check of accept-language header
func workflow1D_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow1D_Test: Check Accept-Language header")

	// Reset database
	resetDatabase()

	// Step1: UserID1 login
	headers := map[string]string{
		"Accept-Language": "zh",
	}

	// Trying to login with empty mobile and password and Accept-Language: zh header,
	// the response should contain non-ascii characxters
	{
		_, err := mobileCredentialsLoginWithHeaders(r, "", "", &headers)
		tests.AssertNotNil(t, err)
		errMessage := err.Error()
		seenNonAscii := false
		for _, c := range errMessage {
			if c > unicode.MaxASCII {
				seenNonAscii = true
				break
			}
		}
		tests.AssertEqual(t, seenNonAscii, true)
	}

	logWorkflowSuccess("Workflow1D_Test ended")
}
