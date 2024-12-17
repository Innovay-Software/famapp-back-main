package integrationTests

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app"
	"github.com/innovay-software/famapp-main/tests"
)

// Workflow5 tests for admin tasks
func TestWorkflow5(t *testing.T) {
	if !runWorkflow5 {
		return
	}

	_, b, _, _ := runtime.Caller(0)
	projDir := filepath.Dir(b)
	r, _ := app.InitApiIntegrationTestServer(fmt.Sprintf("%s/../..", projDir))

	workflow5B_Test(t, r)
	workflow5C_Test(t, r)

	// Reset database
	resetDatabase()
}

// Workflow5B: Admin login and manage members
func workflow5B_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow5B_Test: Login as admin and manage members")

	// Reset database
	resetDatabase()

	// UserID1 login
	res1, _ := mobileCredentialsLogin(r, superAdminMobile, superAdminPassword)
	adminToken := res1.AccessToken

	res1B, err1 := adminListMembersAPI(r, adminToken)
	tests.AssertNil(t, err1)
	tests.AssertEqual(t, len(res1B.Users), 1)
	tests.AssertEqual(t, (res1B.Users)[0].ID, superAdminID)

	password := "123456"
	passcode := "123456"
	// Add a new member
	res2, err2 := adminAddMemberAPI(r, adminToken, "Test1", "12345678901", "member", password, passcode)
	tests.AssertNil(t, err2)
	tests.AssertEqual(t, res2.Member.Name, "Test1")
	tests.AssertEqual(t, res2.Member.Mobile, "12345678901")
	tests.AssertEqual(t, res2.Member.Role, "member")

	// Add the same member should fail
	_, err := adminAddMemberAPI(r, adminToken, "Test1", "12345678901", "member", password, passcode)
	tests.AssertNotNil(t, err)

	// Login as new member should succeed
	res3, _ := mobileCredentialsLogin(r, "12345678901", "123456")
	tests.AssertNotEqual(t, res3.AccessToken, "")
	tests.AssertEqual(t, res3.User.Name, "Test1")
	tests.AssertEqual(t, res3.User.Mobile, "12345678901")

	// Update Member
	res4, _ := adminUpdateMemberAPI(r, adminToken, res3.User.ID, "Test2", "", "member", nil, nil)
	tests.AssertEqual(t, res4.Member.Name, "Test2")
	tests.AssertEqual(t, res4.Member.Mobile, "12345678901") // mobile should not be updated

	// Save invalid password should fail
	newpassword := "123"
	res4b, err := adminUpdateMemberAPI(r, adminToken, res3.User.ID, "", "", "", &newpassword, nil)
	tests.AssertNil(t, res4b.Member)
	tests.AssertNotNil(t, err)

	// Delete Member
	adminDeleteMemberAPI(r, adminToken, res4.Member.UUID.String())

	// Check if it is deleted
	res5, _ := adminListMembersAPI(r, adminToken)
	tests.AssertEqual(t, len(res5.Users), 1)
	tests.AssertEqual(t, (res5.Users)[0].ID, superAdminID)

	logWorkflowSuccess("Workflow5B_Test ended")
}

// Workflow5C: Check for admin permissions
func workflow5C_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow5C_Test: Check user admin permission")

	// Reset database
	resetDatabase()

	// Step1: UserID1 login
	res1, _ := mobileCredentialsLogin(r, superAdminMobile, superAdminPassword)
	superAdminToken := res1.AccessToken

	password := "123456"
	passcode := "123456"

	// Step2: Add member user
	res2, _ := adminAddMemberAPI(r, superAdminToken, "Member1", "12341", "member", password, passcode)
	res2b, _ := mobileCredentialsLogin(r, "12341", "123456")
	member1Token := res2b.AccessToken

	// Step3: Add admin user
	res3, _ := adminAddMemberAPI(r, superAdminToken, "Admin1", "12342", "admin", password, passcode)
	tests.AssertNotNil(t, res3)
	res3b, _ := mobileCredentialsLogin(r, "12342", "123456")
	admin1Token := res3b.AccessToken
	tests.AssertNotEqual(t, admin1Token, "")

	// Step4: member call admin functions, should fail because member is not an admin
	_, err4a := adminListMembersAPI(r, member1Token)
	tests.AssertNotNil(t, err4a)
	_, err4b := adminAddMemberAPI(r, member1Token, "Member2", "12343", "member", password, passcode)
	tests.AssertNotNil(t, err4b)
	_, err4c := adminDeleteMemberAPI(r, member1Token, res2.Member.UUID.String())
	tests.AssertNotNil(t, err4c)

	// Step5: admin call admin functions
	res5a, err5a := adminListMembersAPI(r, admin1Token)
	tests.AssertNil(t, err5a)
	tests.AssertEqual(t, len(res5a.Users), 3)
	res5b, err5b := adminAddMemberAPI(r, admin1Token, "Member2", "12343", "member", password, passcode)
	tests.AssertNil(t, err5b)
	_, err5c := adminDeleteMemberAPI(r, admin1Token, res5b.Member.UUID.String())
	tests.AssertNil(t, err5c)

	// Step6: ther should be 4 users in total
	res6a, err6a := adminListMembersAPI(r, superAdminToken)
	tests.AssertNil(t, err6a)
	tests.AssertNotNil(t, res6a)
	tests.AssertEqual(t, len(res6a.Users), 3)

	logWorkflowSuccess("Workflow5C_Test ended")
}
