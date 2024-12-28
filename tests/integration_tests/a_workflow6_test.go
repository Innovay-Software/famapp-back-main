package integrationTests

import (
	"fmt"
	// "os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/tests"
)

// Workflow6 - Utils related tests
func TestWorkflow6(t *testing.T) {
	if !runWorkflow6 {
		return
	}

	_, b, _, _ := runtime.Caller(0)
	projDir := filepath.Dir(b)
	r, _ := app.InitApiIntegrationTestServer(fmt.Sprintf("%s/../..", projDir))

	workflow6A_Test(t, r)
	workflow6B_Test(t, r)
	workflow6C_Test(t, r)
	workflow6D_Test(t, r)

	// Reset database
	resetDatabase()
}

// Workflow6A: Configs
func workflow6A_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow6A_Test: Test get config API")

	// Reset database
	resetDatabase()

	errC := insertConfig(&models.Config{
		ConfigKey:   "C1",
		ConfigValue: "V1",
		ConfigType:  "text",
	})
	tests.AssertNil(t, errC)

	res1, _ := mobileCredentialsLogin(r, superAdminMobile, superAdminPassword)
	token := res1.AccessToken

	// Get a non existing config and query it without log in - expects permission denied
	_, err2 := utilGetConfigAPI(r, "", "abcde")
	tests.AssertNotNil(t, err2)

	// Get a existing config and query it without log in - expects permission denied
	_, err3 := utilGetConfigAPI(r, "", "C1")
	tests.AssertNotNil(t, err3)

	// Get a non existing config and query it after log in - expects empty value
	res4, err4 := utilGetConfigAPI(r, token, "DFDSREW")
	tests.AssertNil(t, err4)
	tests.AssertEqual(t, res4.ConfigValue, "")

	// Get a existing config and query it after log in - expects config value
	res5, err5 := utilGetConfigAPI(r, token, "C1")
	tests.AssertNil(t, err5)
	tests.AssertEqual(t, res5.ConfigValue, "V1")

	logWorkflowSuccess("Workflow6A_Test ended")
}

// Workflow6B: Ping
func workflow6B_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow6B_Test: Test ping API")

	// Reset database
	resetDatabase()

	res, err := utilPingAPI(r)
	tests.AssertNil(t, err)
	tests.AssertEqual(t, strings.ToLower(res.Ping), "pong")

	logWorkflowSuccess("Workflow6B_Test ended")
}

// Workflow6C: Check for update
func workflow6C_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow6C_Test: Test ping API")

	// Reset database
	resetDatabase()

	// No app version in dastabase - expects no updates
	res1, err1 := utilCheckForMobileUpdate(r, "ios", "1")
	tests.AssertNil(t, err1)
	tests.AssertEqual(t, res1.HasUpdate, false)
	tests.AssertEqual(t, res1.ForceUpdate, false)

	// Insert ios version 3.6.8, effective now
	iosAppVersion := models.AppVersion{
		OS:          "ios",
		Version:     "3.6.8",
		Title:       "3.6.8 Released",
		Content:     map[string]any{"note:": "3.6.8 Released Content"},
		Link:        "https://download.com/latest",
		EffectiveOn: time.Now(),
	}

	// Insert android version 3.6.8, effective now
	androidAppVersion1 := models.AppVersion{
		OS:          "android",
		Version:     "3.6.8",
		Title:       "3.6.8 Released",
		Content:     map[string]any{"note:": "3.6.8 Released Content"},
		Link:        "https://download.com/latest",
		EffectiveOn: time.Now(),
	}

	// Insert android version 4.1.0, effective 24 hours from now
	androidAppVersion2 := models.AppVersion{
		OS:          "android",
		Version:     "4.1.0",
		Title:       "4.1.0 Released",
		Content:     map[string]any{"note:": "4.1.0 Released Content"},
		Link:        "https://download.com/latest",
		EffectiveOn: time.Now().Add(time.Second * 3600 * 24),
	}

	errA := insertAppVersion(&iosAppVersion)
	errB := insertAppVersion(&androidAppVersion1)
	errC := insertAppVersion(&androidAppVersion2)
	tests.AssertNil(t, errA)
	tests.AssertNil(t, errB)
	tests.AssertNil(t, errC)

	// Same version - expects no updates at all
	{
		res, err := utilCheckForMobileUpdate(r, "ios", "3.6.8")
		tests.AssertNil(t, err)
		tests.AssertEqual(t, res.HasUpdate, false)
		tests.AssertEqual(t, res.ForceUpdate, false)
	}

	// Minor version update - expects voluntary update
	{
		res, err := utilCheckForMobileUpdate(r, "ios", "3.5.8")
		tests.AssertNil(t, err)
		tests.AssertEqual(t, res.HasUpdate, true)
		tests.AssertEqual(t, res.ForceUpdate, false)
	}

	// Major version update - expects force update
	{
		res, err := utilCheckForMobileUpdate(r, "ios", "2.9.8")
		tests.AssertNil(t, err)
		tests.AssertEqual(t, res.HasUpdate, true)
		tests.AssertEqual(t, res.ForceUpdate, true)
		tests.AssertEqual(t, res.Url, "https://download.com/latest")
	}

	// Same version - expects no updates at all
	{
		res, err := utilCheckForMobileUpdate(r, "android", "2.6.8")
		tests.AssertNil(t, err)
		tests.AssertEqual(t, res.HasUpdate, true)
		tests.AssertEqual(t, res.ForceUpdate, true)
		tests.AssertEqual(t, res.Version, "3.6.8")
	}

	// Higher version - expects no update because higher version is not effective
	{
		res, err := utilCheckForMobileUpdate(r, "android", "3.6.8")
		tests.AssertNil(t, err)
		tests.AssertEqual(t, res.HasUpdate, false)
		tests.AssertEqual(t, res.ForceUpdate, false)
	}

	// Higher version - expects no update because higher version is not effective
	{
		res, err := utilCheckForMobileUpdate(r, "android", "4.1.8")
		tests.AssertNil(t, err)
		tests.AssertEqual(t, res.Version, "3.6.8")
		tests.AssertEqual(t, res.HasUpdate, false)
		tests.AssertEqual(t, res.ForceUpdate, false)
	}

	logWorkflowSuccess("Workflow6C_Test ended")
}

// Workflow6C: Check for util base64 chunked upload
func workflow6D_Test(t *testing.T, r *gin.Engine) {
	logWorkflowSuccess("Workflow6D_Test: Utils/base64-chunked-upload")

	// Reset database
	resetDatabase()

	// User login
	loginRes1, _ := mobileCredentialsLogin(r, superAdminMobile, superAdminPassword)
	user1AccessToken1 := loginRes1.AccessToken

	// Upload image file
	{
		uploadedFilePath, fileUrl, err := utilBase64ChunkUploadFileWrapper(t, r, user1AccessToken1, "../files/sample-image.jpeg")
		tests.AssertNil(t, err)
		tests.AssertNotEqual(t, uploadedFilePath, "")
		tests.AssertNotEqual(t, fileUrl, "")
		// defer os.Remove(uploadedFilePath)
	}

	// Upload video file
	{
		uploadedFilePath, fileUrl, err := utilBase64ChunkUploadFileWrapper(t, r, user1AccessToken1, "../files/sample-video.mp4")
		tests.AssertNil(t, err)
		tests.AssertNotEqual(t, uploadedFilePath, "")
		tests.AssertNotEqual(t, fileUrl, "")
		// defer os.Remove(uploadedFilePath)
	}

	logWorkflowSuccess("Workflow6D_Test ended")
}
