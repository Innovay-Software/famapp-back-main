package unitTests

import (
	"fmt"
	"os"
	"testing"

	"github.com/innovay-software/famapp-main/app/services"
	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/innovay-software/famapp-main/config"
	"github.com/innovay-software/famapp-main/tests"
)

// Workflow1 tests for user authentication, user profile
func TestFolderFileService(t *testing.T) {
	if !runFolderFileServiceTests {
		return
	}

	// _, b, _, _ := runtime.Caller(0)
	// projDir := filepath.Dir(b)
	// r, _ := app.InitApiIntegrationTestServer(fmt.Sprintf("%s/../..", projDir))

	extractImageExifData_Test(t, "../files/sample-image.jpeg")
	extractVideoExifData_Test(t, "../files/sample-video.mp4")
}

// Extract Image Data
func extractImageExifData_Test(t *testing.T, imageFilepath string) {
	logWorkflowSuccess("ExtractImageExifData_Test: Started")

	// Test basic EXIF extract
	{
		res := services.ExtractFileMetadata(imageFilepath)
		tests.AssertNotNil(t, *res)
		tests.AssertEqual(t, (*res)["dimension"], "1920x1080")
		tests.AssertEqual(t, fmt.Sprintf("%v", (*res)["size"]), "14679474")
		utils.LogSuccess("Basic EXIF extract passed")
	}

	// Compress Image File
	{
		// make a duplicate
		compressedFilePath := imageFilepath + ".compressed.jpg"
		err := utils.DuplicateFile(imageFilepath, compressedFilePath)
		tests.AssertNil(t, err)
		defer os.Remove(compressedFilePath)

		_, err2 := utils.CompressImageToJpgWithMaxSize(compressedFilePath, 960)
		tests.AssertNil(t, err2)

		w, h := utils.GetImageDimeision(compressedFilePath)
		tests.AssertEqual(t, w, 960)
		tests.AssertEqual(t, h, 431)
		utils.LogSuccess("Compress Image File passed")
	}

	// Generate Thumbnail Image File
	{
		thumbnailFilepath := imageFilepath + ".thumbnail.jpg"
		_, err := utils.GenerateThumbnailJpg(imageFilepath, thumbnailFilepath)
		tests.AssertNil(t, err)
		defer os.Remove(thumbnailFilepath)

		// Check thumbnail dimension
		w, h := utils.GetImageDimeision(thumbnailFilepath)
		tests.AssertEqual(t, w, config.FolderFileThumbnailSize)
		tests.AssertEqual(t, h, config.FolderFileThumbnailSize)
		utils.LogSuccess("Thumbnail Image File passed")
	}

	logWorkflowSuccess("ExtractImageExifData_Test ended")
}

// Extract Video Data
func extractVideoExifData_Test(t *testing.T, videoFilepath string) {
	logWorkflowSuccess("ExtractVideoExifData_Test: Started")

	{
		res := services.ExtractFileMetadata(videoFilepath)
		tests.AssertNotNil(t, *res)
		tests.AssertEqual(t, (*res)["dimension"], "1920x1080")
		tests.AssertEqual(t, fmt.Sprintf("%v", (*res)["size"]), "17839845")

		duration, err := utils.ExtractVideoDuration(videoFilepath)
		tests.AssertNil(t, err)
		tests.AssertEqual(t, duration, 31)
	}

	// Compresse video file
	{
		compressedFilePath := "../files/sample.mp4.compressed.mp4"
		err := utils.DuplicateFile(videoFilepath, compressedFilePath)
		tests.AssertNil(t, err)
		defer os.Remove(compressedFilePath)

		_, err2 := utils.CompressVideoToMp4FullHD(compressedFilePath)
		tests.AssertNil(t, err2)
	}

	// Generate thumbnail
	{
		firstFrameFilepath := videoFilepath + ".frame.jpg"
		thumbnailFilepath := videoFilepath + ".thumbnail.jpg"
		_, err := utils.ExtractVideoFirstFrameAsJpg(videoFilepath, firstFrameFilepath)
		tests.AssertNil(t, err)
		defer os.Remove(firstFrameFilepath)

		w1, h1 := utils.GetImageDimeision(firstFrameFilepath)
		tests.AssertEqual(t, w1, 1920)
		tests.AssertEqual(t, h1, 1080)

		_, err2 := utils.GenerateThumbnailJpg(firstFrameFilepath, thumbnailFilepath)
		tests.AssertNil(t, err2)
		defer os.Remove(thumbnailFilepath)

		w2, h2 := utils.GetImageDimeision(thumbnailFilepath)
		tests.AssertEqual(t, w2, config.FolderFileThumbnailSize)
		tests.AssertEqual(t, h2, config.FolderFileThumbnailSize)
	}

	logWorkflowSuccess("ExtractVideoExifData_Test ended")
}
