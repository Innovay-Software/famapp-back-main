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

// Workflow3 Tests for FolderFiles
func TestWorkflow3(t *testing.T) {
	if !runWorkflow3 {
		return
	}

	_, b, _, _ := runtime.Caller(0)
	projDir := filepath.Dir(b)
	r, _ := app.InitApiIntegrationTestServer(fmt.Sprintf("%s/../..", projDir))

	Workflow3A_Test(t, r)
}

// Folder - Basic CRUD
func Workflow3A_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow3A_Test: Folder - Basic")

	// Reset database
	resetDatabase()

	// Login as admin
	user1Id := uint64(0)
	user1Token := ""
	adminToken := ""
	{
		res, _ := mobileCredentialsLogin(r, superAdminMobile, superAdminPassword)
		adminToken = res.AccessToken
	}

	// Create new user
	{
		res1, err1 := createUser(r, "User1", "1234561", "member")
		tests.AssertNil(t, err1)
		tests.AssertEqual(t, res1.Member.Mobile, "1234561")
	}

	// Login as the new user
	{
		res1b, err1b := mobileCredentialsLogin(r, "1234561", "123456")
		tests.AssertNil(t, err1b)
		user1Token = res1b.AccessToken
		user1Id = res1b.User.ID
		tests.AssertNotEqual(t, user1Token, "")
		folders := res1b.Folders
		tests.AssertEqual(t, len(folders), 0)
	}

	// Add folder
	folderID := uint64(0)
	{
		res2, err2 := saveFolder(
			r, user1Token, 0, user1Id, 0, "Test1", "", "normal",
			false, false, &map[string]any{}, []uint64{})
		// folderID := (*res2.Folder).ID
		tests.AssertNil(t, err2)
		tests.AssertNotEqual(t, (*res2.Folder).ID, 0)
		tests.AssertEqual(t, (*res2.Folder).Title, "Test1")
		folderID = (*res2.Folder).ID
	}

	// Update folder
	{
		res3, err3 := saveFolder(
			r, user1Token, folderID, user1Id, 0, "Test2", "", "normal",
			true, true, &map[string]any{}, []uint64{})
		tests.AssertNil(t, err3)
		tests.AssertEqual(t, (*res3.Folder).Title, "Test2")
	}

	// Create User2
	user2Id := uint64(0)
	user2Token := ""
	{
		res4, err4 := createUser(r, "User2", "1234562", "member")
		tests.AssertNil(t, err4)
		tests.AssertEqual(t, res4.Member.Mobile, "1234562")
		user2Id = res4.Member.ID
	}

	// Login as the new user
	{
		res4b, err4b := mobileCredentialsLogin(r, "1234562", "123456")
		tests.AssertNil(t, err4b)
		user2Token = res4b.AccessToken
		tests.AssertNotEqual(t, user2Token, "")
	}

	// Save folder with user2 should fail
	{
		_, err4c := saveFolder(
			r, user2Token, folderID, user1Id, 0, "Test2", "", "normal",
			true, true, &map[string]any{}, []uint64{user2Id})
		tests.AssertNotNil(t, err4c)
	}

	// User1 share folder with user2
	{
		res5, err5 := saveFolder(
			r, user1Token, folderID, user1Id, 0, "Test2", "", "normal",
			true, true, &map[string]any{}, []uint64{user2Id})
		tests.AssertNil(t, err5)
		tests.AssertEqual(t, (*res5.Folder).Title, "Test2")
	}

	// Save folder with user2 should fail again,
	// because invitees only have view permission not update permission
	{
		_, err5b := saveFolder(
			r, user2Token, folderID, user1Id, 0, "Test2222", "", "normal",
			true, true, &map[string]any{}, []uint64{})
		tests.AssertNotNil(t, err5b)
	}

	{
		res6, err6 := saveFolder(
			r, adminToken, folderID, user1Id, 0, "Test2222", "", "normal",
			true, true, &map[string]any{}, []uint64{})
		tests.AssertNil(t, err6)
		tests.AssertEqual(t, res6.Folder.Title, "Test2222")
	}

	logWorkflowSuccess("Workflow3A_Test: ended")
}
