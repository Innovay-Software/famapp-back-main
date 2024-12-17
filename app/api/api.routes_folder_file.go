/**
ApiServerInterfaceImpl - Folder and FolderFiles Part
Handles all Folder and FolderFiles related routes
*/

package api

import (
	"github.com/gin-gonic/gin"
	folderFileHandlers "github.com/innovay-software/famapp-main/app/handlers/folderFile"
)

// Delete FolderFiles
func (s *ApiServerInterfaceImpl) FolderFileDeleteFolderFilesPath(
	c *gin.Context, folderId int64, params FolderFileDeleteFolderFilesPathParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	var req FolderFileDeleteFolderFilesPathJSONRequestBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := folderFileHandlers.DeleteFolderFilesHandler(
		c, user, folderId, req.FolderFileIds,
	)
	handleApiResponse(c, res, err)
}

// Update multiple FolderFiles
func (s *ApiServerInterfaceImpl) FolderFileUpdateMultipleFolderFilesPath(
	c *gin.Context, params FolderFileUpdateMultipleFolderFilesPathParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	var req FolderFileUpdateMultipleFolderFilesPathJSONRequestBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := folderFileHandlers.UpdateMultipleFolderFilesHandler(
		c, user, req.FolderFileIds, req.NewFolderId, req.NewShotAtTimestamp,
	)
	handleApiResponse(c, res, err)
}

// Update single FolderFile
func (s *ApiServerInterfaceImpl) FolderFileUpdateSingleFolderFilePath(
	c *gin.Context, params FolderFileUpdateSingleFolderFilePathParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	var req FolderFileUpdateSingleFolderFilePathJSONRequestBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := folderFileHandlers.UpdateSingleFolderFileHandler(
		c, user, req.FolderFileId, req.IsPrivate, req.Remark,
	)
	handleApiResponse(c, res, err)
}

// Ger FolderFiles after target datetime
func (s *ApiServerInterfaceImpl) FolderFileGetFolderFilesAfterMicroTimestampShotAt(
	c *gin.Context, folderId int64, pivotDate string, microtimestamp int64,
	params FolderFileGetFolderFilesAfterMicroTimestampShotAtParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request body

	// Handle request
	res, err := folderFileHandlers.GetFolderFilesAfterShotAtHandler(
		c, user, folderId, pivotDate, microtimestamp,
	)
	handleApiResponse(c, res, err)
}

// Get FolderFiles before target datetime
func (s *ApiServerInterfaceImpl) FolderFileGetFolderFilesBeforeMicroTimestampShotAt(
	c *gin.Context, folderId int64, pivotDate string, microtimestamp int64,
	params FolderFileGetFolderFilesBeforeMicroTimestampShotAtParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request body

	// Handle request
	res, err := folderFileHandlers.GetFolderFilesBeforeShotAtHandler(
		c, user, folderId, pivotDate, microtimestamp,
	)
	handleApiResponse(c, res, err)
}

// Delete Folder
func (s *ApiServerInterfaceImpl) FolderFileDeleteFolderPath(
	c *gin.Context, folderId int64, params FolderFileDeleteFolderPathParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request body

	// Handle request
	res, err := folderFileHandlers.DeleteFolderHandler(
		c, user, folderId,
	)
	handleApiResponse(c, res, err)
}

// Save Folder
func (s *ApiServerInterfaceImpl) FolderFileSaveFolderPath(
	c *gin.Context, folderId int64, params FolderFileSaveFolderPathParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	var req FolderFileSaveFolderPathJSONRequestBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := folderFileHandlers.SaveFolderHandler(
		c, user, folderId, req.OwnerId, req.ParentId, req.Title, req.Cover,
		req.Type, req.IsDefault, req.IsPrivate, &req.Metadata, &req.InviteeIds,
	)
	handleApiResponse(c, res, err)
}

// Display FolderFile-Original
func (s *ApiServerInterfaceImpl) FolderFileDisplayOriginalPath(
	c *gin.Context, folderFileId int64, params FolderFileDisplayOriginalPathParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request body

	// Handle request
	filepath, err := folderFileHandlers.DisplayFolderFileOriginalHandler(
		c, user, folderFileId,
	)
	handleFileResponse(c, filepath, err)
}

// Display FolderFile-Compressed(default)
func (s *ApiServerInterfaceImpl) FolderFileDisplayPath(c *gin.Context, folderFileId int64, params FolderFileDisplayPathParams) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request body

	// Handle request
	filepath, err := folderFileHandlers.DisplayFolderFileCompressedHandler(
		c, user, folderFileId,
	)
	handleFileResponse(c, filepath, err)
}

// Display FolderFile-Thumbnail
func (s *ApiServerInterfaceImpl) FolderFileDisplayThumbnailPath(c *gin.Context, folderFileId int64, params FolderFileDisplayThumbnailPathParams) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	// pass - no request body

	// Handle request
	filepath, err := folderFileHandlers.DisplayFolderFileOriginalHandler(
		c, user, folderFileId,
	)
	handleFileResponse(c, filepath, err)
}
