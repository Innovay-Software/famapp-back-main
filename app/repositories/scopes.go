package repositories

import "gorm.io/gorm"

// General
func commonScopeSoftDeleted(query *gorm.DB) *gorm.DB {
	return query.Where("deleted_at IS NOT NULL")
}

// FolderFile scope where folderId > 0
func folderFileScopeActive(query *gorm.DB) *gorm.DB {
	return query.Where("folder_id > ?", 0)
}

// FolderFile scope for ready to upload to google
func folderFileScopeReadyForGoogleUpload(query *gorm.DB) *gorm.DB {
	return query.Where("hw_original_file_path IS NOT NULL").
		Where("hw_original_file_path != ?", "processing").
		Where("is_downloading = ?", 0)
}

// FolderFile scope where is not downloading
func folderFileScopeNotDownloading(query *gorm.DB) *gorm.DB {
	return query.Where("is_downloading = ?", 0)
}

// UserScope where user is admin
func userScopeAdmin(query *gorm.DB) *gorm.DB {
	return query.Where("role = ?", "admin")
}

// UserScope where user is manager
func userScopeManager(query *gorm.DB) *gorm.DB {
	return query.Where("role = ?", "manager")
}

// UserScope where user is either admin or manager
func userScopeAdminOrManager(query *gorm.DB) *gorm.DB {
	return query.Where("role in ?", []string{"admin", "manager"})
}
