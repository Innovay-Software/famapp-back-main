package routers

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/innovay-software/famapp-main/app/handlers"
// )

// type RouteNode struct {
// 	Val      string
// 	Children *[]RouteNode
// }

// // Register API routes
// func registerApiV2Routes(router *gin.Engine) {
// 	v2 := router.Group("/api/v2")
// 	// uploadGroup := v2.Group("/bin/upload")
// 	// {
// 	// 	uploadGroup.POST("/chunk-upload-folder-file/:folderId", userHandlerWrapper(handlers.FolderFileChunkUploadFileHandler))
// 	// }
// 	// oauthGroup := v2.Group("/auth")
// 	// {
// 	// 	oauthGroup.POST("/mobile-login", handlerWrapper(handlers.LoginHandler))
// 	// 	oauthGroup.POST("/access-token-login", userHandlerWrapper(handlers.AccessTokenLoginHandler))
// 	// 	oauthGroup.POST("/update-profile", userHandlerWrapper(handlers.UpdateUserProfileHandler))
// 	// }
// 	commGroup := v2.Group("/communicator")
// 	{
// 		commGroup.POST("/download-folder-file/:filename/:type", handlerWrapper(handlers.UnimplementedHandler))
// 		commGroup.POST("/sync-folder-file-to-peer-server", handlerWrapper(handlers.UnimplementedHandler))
// 	}
// 	// adminGroup := v2.Group("/admin")
// 	// {
// 	// 	adminGroup.POST("/get-all-users/:afterId", adminHandlerWrapper(handlers.AdminGetMemberListHandler))
// 	// 	adminGroup.POST("/save-user/:userId", adminHandlerWrapper(handlers.AdminSaveMember))
// 	// 	adminGroup.POST("/delete-user/:userId", adminHandlerWrapper(handlers.AdminDeleteMember))
// 	// }
// 	userGroup := v2.Group("/user")
// 	{
// 		userGroup.POST("/set-user-avatar", handlerWrapper(handlers.UnimplementedHandler))
// 		userGroup.POST("/get-users", handlerWrapper(handlers.UnimplementedHandler))
// 		userGroup.POST("/update-device-token", handlerWrapper(handlers.UnimplementedHandler))
// 	}
// 	utilGroup := v2.Group("/util")
// 	{
// 		utilGroup.GET("/ping", handlerWrapper(handlers.Ping))
// 		utilGroup.POST("/ping", handlerWrapper(handlers.Ping))
// 		utilGroup.GET("/user-avatar/:userId", fileHandlerWrapper(handlers.UserAvatar))
// 		utilGroup.GET("/check-for-update/:os/:currentVersion", handlerWrapper(handlers.CheckForUpdate))
// 		utilGroup.POST("/check-for-update/:os/:currentVersion", handlerWrapper(handlers.CheckForUpdate))

// 		utilGroup.POST("/base64-chunked-upload-file", userHandlerWrapper(handlers.Base64ChunkUploadFileHandler))
// 		utilGroup.POST("/upload-image", handlerWrapper(handlers.UnimplementedHandler))
// 		utilGroup.POST("/upload-file", handlerWrapper(handlers.UnimplementedHandler))
// 		utilGroup.POST("/config/:configKey", handlerWrapper(handlers.UnimplementedHandler))
// 	}
// 	// FolderFiles Group
// 	ffGroup := v2.Group("/files")
// 	{
// 		// ffGroup.POST("/get-folder-files-before-id/:folderId/:pivotDate", userHandlerWrapper(handlers.GetFolderFilesBeforeIdHandler))
// 		// ffGroup.POST("/get-folder-files-after-id/:folderId/:pivotDate", userHandlerWrapper(handlers.GetFolderFilesAfterIdHandler))
// 		ffGroup.POST("/save-folder/:folderId", userHandlerWrapper(handlers.SaveFolderHandler))
// 		ffGroup.POST("/delete-folder/:folderId", userHandlerWrapper(handlers.DeleteFolderHandler))
// 		// ffGroup.POST("/delete-files/:folderId", userHandlerWrapper(handlers.DeleteFolderFilesHandler))
// 		// ffGroup.POST("/update-folder-file/:folderFileId", userHandlerWrapper(handlers.UpdateSingleFolderFileHandler))
// 		// ffGroup.POST("/update-multiple-folder-files", userHandlerWrapper(handlers.UpdateMultipleFolderFilesHandler))

// 		// ffGroup.POST("/move-files-to-folder/:folderId", userHandlerWrapper(handlers.UnimplementedHandler))
// 		// ffGroup.POST("/set-shot-at-date/:date", userHandlerWrapper(handlers.UnimplementedHandler))
// 	}
// 	// Locker Group
// 	// lGroup := v2.Group("/locker-notes")
// 	// {
// 	// 	lGroup.POST("/list-notes", userHandlerWrapper(handlers.ListLockerNotesHandler))
// 	// 	lGroup.POST("/save-note/:noteId", userHandlerWrapper(handlers.SaveLockerNoteHandler))
// 	// 	lGroup.POST("/delete-note/:noteId", userHandlerWrapper(handlers.DeleteLockerNoteHandler))
// 	// }
// 	// FolderFileDisplay group
// 	ffdGroup := v2.Group("/folder-file-display")
// 	{
// 		ffdGroup.GET("/folder-file/:folderFileId", userFileHandlerWrapper(handlers.DisplayFolderFileCompressedHandler))
// 		ffdGroup.GET("/folder-file-thumbnail/:folderFileId", userFileHandlerWrapper(handlers.DisplayFolderFileThumbnailHandler))
// 		ffdGroup.GET("/folder-file-original/:folderFileId", userFileHandlerWrapper(handlers.DisplayFolderFileOriginalHandler))
// 	}
// }

// func handlerWrapper(handler handlers.ApiHandler) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		handlers.ApiHandlerWrapper(c, handler)
// 	}
// }

// func fileHandlerWrapper(handler handlers.ApiFileHandler) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		handlers.ApiFileHandlerWrapper(c, handler)
// 	}
// }

// func userHandlerWrapper(handler handlers.ApiUserHandler) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		handlers.ApiUserHandlerWrapper(c, handler)
// 	}
// }

// func adminHandlerWrapper(handler handlers.ApiAdminHandler) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		handlers.ApiAdminHandlerWrapper(c, handler)
// 	}
// }

// func userFileHandlerWrapper(handler handlers.ApiUserFileHandler) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		handlers.ApiUserFileHandlerWrapper(c, handler)
// 	}
// }
