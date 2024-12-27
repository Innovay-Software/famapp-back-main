/**
ApiServerInterfaceImpl - Folder and FolderFiles Part
Handles all Folder and FolderFiles related routes
*/

package api

import (
	"github.com/gin-gonic/gin"
	folderFileHandlers "github.com/innovay-software/famapp-main/app/handlers/folderFile"
	"github.com/innovay-software/famapp-main/app/utils"
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

	folderFileIds := utils.Int64SliceToUint64(&(req.FolderFileIds))
	// Handle request
	res, err := folderFileHandlers.DeleteFolderFilesHandler(
		c, user, uint64(folderId), *folderFileIds,
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

	folderFileIds := utils.Int64SliceToUint64(&(req.FolderFileIds))
	newFolderId := utils.Int64PointerToUint64Pointer(req.NewFolderId)

	// Handle request
	res, err := folderFileHandlers.UpdateMultipleFolderFilesHandler(
		c, user, *folderFileIds, newFolderId, req.NewTakenOnTimestamp,
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
		c, user, uint64(req.FolderFileId), req.IsPrivate, req.Remark,
	)
	handleApiResponse(c, res, err)
}

// Ger FolderFiles after target datetime
func (s *ApiServerInterfaceImpl) FolderFileGetFolderFilesAfterMicroTimestampTakenOn(
	c *gin.Context, folderId int64, pivotDate string, microtimestamp int64,
	params FolderFileGetFolderFilesAfterMicroTimestampTakenOnParams,
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
	res, err := folderFileHandlers.GetFolderFilesAfterTakenOnHandler(
		c, user, uint64(folderId), pivotDate, microtimestamp,
	)
	handleApiResponse(c, res, err)
}

// Get FolderFiles before target datetime
func (s *ApiServerInterfaceImpl) FolderFileGetFolderFilesBeforeMicroTimestampTakenOn(
	c *gin.Context, folderId int64, pivotDate string, microtimestamp int64,
	params FolderFileGetFolderFilesBeforeMicroTimestampTakenOnParams,
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
	res, err := folderFileHandlers.GetFolderFilesBeforeTakenOnHandler(
		c, user, uint64(folderId), pivotDate, microtimestamp,
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
		c, user, uint64(folderId),
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

	inviteeIds := utils.Int64SliceToUint64(&req.InviteeIds)

	// Handle request
	res, err := folderFileHandlers.SaveFolderHandler(
		c, user, uint64(folderId), uint64(req.OwnerId), uint64(req.ParentId), req.Title, req.Cover,
		req.Type, req.IsDefault, req.IsPrivate, &req.Metadata, inviteeIds,
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
		c, user, uint64(folderFileId),
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
		c, user, uint64(folderFileId),
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
		c, user, uint64(folderFileId),
	)
	handleFileResponse(c, filepath, err)
}
