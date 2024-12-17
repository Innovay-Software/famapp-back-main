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

// Workflow2 tests for LockerNotes
func TestWorkflow2(t *testing.T) {
	if !runWorkflow2 {
		return
	}

	_, b, _, _ := runtime.Caller(0)
	projDir := filepath.Dir(b)
	r, _ := app.InitApiIntegrationTestServer(fmt.Sprintf("%s/../..", projDir))

	workflow2A_Test(t, r)
	workflow2B_Test(t, r)
}

// LockerNotes - Basic CRUD
func workflow2A_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("workflow2A_Test: Locker Notes - Basic")

	// Reset database
	resetDatabase()

	// Create new user - expects success
	{
		res1, err1a := createUser(r, "User1", "1234561", "member")
		tests.AssertNil(t, err1a)
		tests.AssertEqual(t, res1.Member.Mobile, "1234561")
	}

	// Login as the new user - expects success
	user1Token := ""
	{
		res1b, err1b := mobileCredentialsLogin(r, "1234561", "123456")
		tests.AssertNil(t, err1b)
		user1Token = res1b.AccessToken
		tests.AssertNotEqual(t, user1Token, "")
	}

	// List locker notes - expects empty
	{
		res2, err2 := listLockerNotes(r, user1Token)
		tests.AssertNil(t, err2)
		tests.AssertEqual(t, len(*res2.Notes), 0)
	}

	// Create a new note - expects success
	noteId := int64(0)
	{
		res3, err3 := createLockerNotes(r, user1Token, "Title1", "Content1", []int64{})
		tests.AssertNil(t, err3)
		tests.AssertNotNil(t, res3.Note)
		tests.AssertEqual(t, res3.Note.Title, "Title1")
		noteId = (*res3.Note).ID
	}

	// Update note - expects success
	{
		res4, err4 := saveLockerNotes(r, user1Token, noteId, "Title2", "", []int64{})
		tests.AssertNil(t, err4)
		tests.AssertNotNil(t, res4.Note)
		tests.AssertEqual(t, res4.Note.Title, "Title2")
	}

	// Query all notes - expects 1 note
	{
		res5, err5 := listLockerNotes(r, user1Token)
		tests.AssertNil(t, err5)
		tests.AssertEqual(t, len(*res5.Notes), 1)
		tests.AssertEqual(t, (*res5.Notes)[0].Title, "Title2")
	}

	// Delete note - expects ok
	{
		_, err6 := deleteLockerNotes(r, user1Token, noteId)
		tests.AssertNil(t, err6)
	}

	// Query all notes - expects empty
	{
		res7, err7 := listLockerNotes(r, user1Token)
		tests.AssertNil(t, err7)
		tests.AssertEqual(t, len(*res7.Notes), 0)
	}

	logWorkflowSuccess("workflow2A_Test: ended")
}

// LockerNotes - Sharing between users
func workflow2B_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("workflow2A_Test: Locker Notes - Sharing between users")

	// Reset database
	resetDatabase()

	// Create user1
	user1Token := ""
	{
		res1, err1 := createUser(r, "User1", "1234561", "member")
		tests.AssertNil(t, err1)
		tests.AssertEqual(t, res1.Member.Mobile, "1234561")
	}

	// Login as user1
	{
		res1b, err1b := mobileCredentialsLogin(r, "1234561", "123456")
		tests.AssertNil(t, err1b)
		user1Token = res1b.AccessToken
		tests.AssertNotEqual(t, res1b.AccessToken, "")
		tests.AssertNotEqual(t, res1b.User.ID, 0)
	}

	// Create user2
	user2Id := int64(0)
	user2Token := ""
	{
		res2, err2 := createUser(r, "User2", "1234562", "member")
		tests.AssertNil(t, err2)
		tests.AssertEqual(t, res2.Member.Mobile, "1234562")
	}

	// Login as user2
	{
		res2b, err2b := mobileCredentialsLogin(r, "1234562", "123456")
		tests.AssertNil(t, err2b)
		user2Id = res2b.User.ID
		user2Token = res2b.AccessToken
		tests.AssertNotEqual(t, res2b.AccessToken, "")
	}

	// User1 add a note and share to user2, user2 should be able to see, but unable to modify
	{
		_, err3 := createLockerNotes(r, user1Token, "Note1", "Content1", []int64{user2Id})
		tests.AssertNil(t, err3)
	}

	// User2 should be able to see the record
	noteId := int64(0)
	{
		res4, err4 := listLockerNotes(r, user2Token)
		tests.AssertNil(t, err4)
		tests.AssertEqual(t, len(*res4.Notes), 1)
		tests.AssertEqual(t, (*res4.Notes)[0].Title, "Note1")
		noteId = (*res4.Notes)[0].ID
	}

	// User2 should not be able to modify
	{
		_, err5 := saveLockerNotes(r, user2Token, noteId, "Note1Updated", "Content1Updated", []int64{user2Id})
		tests.AssertNotNil(t, err5)
	}

	// User1 should be able to modify, and remove user2's access
	{
		_, err6 := saveLockerNotes(r, user1Token, noteId, "Note1Updated", "Content1Updated", []int64{})
		tests.AssertNil(t, err6)
	}

	// User2 should not be able to see it anymore
	{
		res7, err7 := listLockerNotes(r, user2Token)
		tests.AssertNil(t, err7)
		tests.AssertEqual(t, len(*res7.Notes), 0)
	}

	logWorkflowSuccess("workflow2A_Test: ended")
}
