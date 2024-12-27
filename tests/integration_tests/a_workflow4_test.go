package integrationTests

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/tests"
)

// Workflow4 FolderFiles
func TestWorkflow4(t *testing.T) {
	if !runWorkflow4 {
		return
	}

	_, b, _, _ := runtime.Caller(0)
	projDir := filepath.Dir(b)
	fmt.Println("projDir =", projDir)
	gin, _ := app.InitApiIntegrationTestServer(fmt.Sprintf("%s/../..", projDir))

	Workflow4A_Test(t, gin)
	// Workflow4B_Test(t, gin)
}

// Basic Folder CRUD
func Workflow4A_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow4A_Test: Folder - Basic CRUD")

	user1, user1AccessToken, folderId := basicSetupWithLoginAndFolder(t, r, "Folder1")
	tests.AssertNotEqual(t, user1AccessToken, "")
	tests.AssertNotEqual(t, folderId, 0)

	// Create a second user
	var user2 *models.User
	{
		res1, err1 := createUser(r, "User2", "1234562", "member")
		tests.AssertNil(t, err1)
		tests.AssertEqual(t, res1.Member.Mobile, "1234562")
		tests.AssertNotEqual(t, res1.Member.ID, 0)
		user2 = res1.Member
	}

	// Query folders
	{
		res, err := mobileAccessTokenLogin(r, user1AccessToken)
		tests.AssertNil(t, err)
		tests.AssertEqual(t, len(res.Folders), 1)
		tests.AssertEqual(t, (res.Folders)[0].Title, "Folder1")
	}

	// Invite User
	{
		res, err := saveFolder(
			r, user1AccessToken, folderId, user1.ID, 0, "Folder1Updated", "",
			"normal", true, false, nil, []uint64{user2.ID},
		)
		tests.AssertNil(t, err)
		tests.AssertEqual(t, res.Folder.ID, folderId)
	}

	// Query folders again and check for invitees
	{
		res, err := mobileAccessTokenLogin(r, user1AccessToken)
		tests.AssertNil(t, err)
		tests.AssertEqual(t, len(res.Folders), 1)
		invitees := (res.Folders)[0].Invitees
		tests.AssertEqual(t, len(invitees), 2)
		tests.AssertIn(t, invitees[0].Name, []any{"User1", "User2"})
	}

	// Login as Invitee and check for folders
	user2AccessToken := ""
	{
		res, _ := mobileCredentialsLogin(r, "1234562", "123456")
		user2AccessToken = res.AccessToken
		tests.AssertEqual(t, len(res.Folders), 1)
	}

	// User2 delete folder - expected permission denied error
	{
		_, err := deleteFolder(r, user2AccessToken, folderId)
		tests.AssertNotNil(t, err)
	}

	{
		_, err := deleteFolder(r, user1AccessToken, folderId)
		tests.AssertNil(t, err)
	}

	// User1 login again - expects empty folders
	{
		res, _ := mobileCredentialsLogin(r, "1234561", "123456")
		tests.AssertEqual(t, len(res.Folders), 0)
	}
}

// Basic Upload
func Workflow4B_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow4B_Test: FolderFile - BasicUpload")

	// _, b, _, _ := runtime.Caller(0)
	// projDir := filepath.Dir(b)

	_, accessToken, folderId := basicSetupWithLoginAndFolder(t, r, "Folder1")
	tests.AssertNotEqual(t, accessToken, "")
	tests.AssertNotEqual(t, folderId, 0)

	// Upload to folder
	res3, err3 := uploadFileToFolderFile(r, accessToken, folderId, "../files/sample-image.jpeg")
	tests.AssertNil(t, err3)
	fmt.Println(res3)
	// panic(errors.New("haha"))

	// Query latest folder files and the uploaded file should be there

	// Share folder with user2 and he should see it

	// Set folder file to private and user2 should not be able to see

	// Admin should be able to see

	// Set folder file to public and Set a remark to folder file and user2 should be able to see

	// user2 should not be able to update folder file

	logWorkflowSuccess("Workflow4B_Test: ended")
}

func basicSetupWithLoginAndFolder(
	t *testing.T, r *gin.Engine, folderName string,
) (
	*models.User, string, uint64,
) {

	// _, b, _, _ := runtime.Caller(0)
	// projDir := filepath.Dir(b)

	// Reset database
	resetDatabase()

	// Create new user
	res1, err1 := createUser(r, "User1", "1234561", "member")
	tests.AssertNil(t, err1)
	tests.AssertEqual(t, res1.Member.Mobile, "1234561")

	// Login as the new user
	res1b, _ := mobileCredentialsLogin(r, "1234561", "123456")
	user1Token := res1b.AccessToken
	user1ID := res1b.User.ID

	// Add folder
	res, _ := saveFolder(r, user1Token, 0, user1ID, 0, folderName, "", "normal",
		false, false, &map[string]any{}, []uint64{})

	return res1b.User, user1Token, res.Folder.ID
}
