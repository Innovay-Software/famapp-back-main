/**
ApiServerInterfaceImpl - Upload Part
Handles all the routes related to folder file uploads
*/

package api

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	folderFileHandlers "github.com/innovay-software/famapp-main/app/handlers/folderFile"
	utilHandlers "github.com/innovay-software/famapp-main/app/handlers/util"
)

// Get ChunkUploadId
func (s *ApiServerInterfaceImpl) FolderFileGetChunkUploadIdPath(
	c *gin.Context, params FolderFileGetChunkUploadIdPathParams,
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
	res, err := folderFileHandlers.FolderFileChunkUploadInitUploadIdHandler(c, user)
	handleApiResponse(c, res, err)
}

// Base64 ChunkUpload
func (s *ApiServerInterfaceImpl) UtilBase64ChunkUploadPath(
	c *gin.Context, params UtilBase64ChunkUploadPathParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	var req UtilBase64ChunkUploadPathJSONRequestBody
	if err := s.validateInputRequest(c, &req); err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Handle request
	res, err := utilHandlers.Base64ChunkUploadFileHandler(
		c, user, req.Base64EncodedContent, req.FileName,
		req.ChunkedFileName, req.HasMore,
	)
	handleApiResponse(c, res, err)
}

// FolderFile Chunk Upload path
func (s *ApiServerInterfaceImpl) FolderFileChunkUploadPath(
	c *gin.Context, folderId int64, params FolderFileChunkUploadPathParams,
) {
	// Authenticate caller
	user, err := getAuthenticatedUser(c)
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}

	// Validate Request
	chunkIndex, err := strconv.Atoi(params.ChunkIndex)
	if err != nil {
		apiRespFailError(c, errors.New("invalid chunk index"), nil)
		return
	}

	// Handle request
	res, err := folderFileHandlers.FolderFileChunkUploadHandler(
		c, user, folderId, params.UploadId,
		params.HasMore == "1", params.Filename, chunkIndex,
	)
	handleApiResponse(c, res, err)
}
