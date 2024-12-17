package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/innovay-software/famapp-main/app/utils"
	"gorm.io/gorm"
)

type FolderFile struct {
	BaseDbModel

	PeerID                 int64          `gorm:"column:peer_id" json:"-" mapstructure:"peerId"`
	OwnerID                int64          `gorm:"column:owner_id" json:"ownerId" mapstructure:"ownerId"`
	FolderID               int64          `gorm:"column:folder_id" json:"folderId" mapstructure:"folderId"`
	Disk                   string         `gorm:"column:disk" json:"-"`
	FileName               string         `gorm:"column:file_name" json:"fileName"`
	FileType               string         `gorm:"column:file_type" json:"fileType"`
	OriginalFilePath       string         `gorm:"column:original_file_path" json:"originalFilePath"`
	CompressedFilePath     string         `gorm:"column:file_path" json:"filePath"`
	ThumbnailPath          string         `gorm:"column:thumbnail_path" json:"thumbnailPath"`
	HwOriginalFilePath     string         `gorm:"column:hw_original_file_path" json:"hwOriginalFilePath"`
	GoogleOriginalFilePath string         `gorm:"column:google_original_file_path" json:"googleOriginalFilePath"`
	GoogleFilePath         string         `gorm:"column:google_file_path" json:"googleFilePath"`
	GoogleThumbnailPath    string         `gorm:"column:google_thumbnail_path" json:"googleThumbnailPath"`
	GoogleDriveFileID      string         `gorm:"column:google_drive_file_id" json:"googleDriveFileId"`
	Remark                 string         `gorm:"column:remark" json:"remark" mapstructure:"remark"`
	Metadata               map[string]any `gorm:"column:metadata;serializer:json" json:"metadata" mapstructure:"metadata"`
	HasExif                bool           `gorm:"column:has_exif; default:0" json:"hasExif"`
	IsPrivate              bool           `gorm:"column:is_private; default:0" json:"isPrivate" mapstructure:"isPrivate"`
	Views                  int64          `gorm:"column:views; default:0" json:"views"`
	SyncedAt               *time.Time     `gorm:"column:synced_at; null" json:"syncedAt"`
	IsDownloading          bool           `gorm:"column:is_downloading; default:0" json:"-"`
	ShotAt                 time.Time      `gorm:"column:shot_at" json:"shotAtOriginal" mapstructure:"shotAt"`
	ShotAtString           string         `gorm:"-" json:"shotAt"`
	IsPreprocessing        int            `gorm:"-" json:"isPreprocessing"`
	MetadataSimple         map[string]any `gorm:"-" json:"metadataSimple"`
}

func (FolderFile) TableName() string {
	return "folder_files"
}

// AfterSave hook
func (ff *FolderFile) AfterSave(tx *gorm.DB) error {
	// After folder file is saved, update the folder's totalFiles, earliest_shot_at, and latest_shot_at

	if ff.FolderID > 0 {
		var folder Folder
		if err := tx.First(&folder, ff.FolderID).Error; err == nil {
			// if file shots at before folder's earliest shot at date, update it
			if folder.EarliestShotAt == nil || ff.ShotAt.Compare(*folder.EarliestShotAt) < 0 {
				folder.EarliestShotAt = &ff.ShotAt
			}
			// if file shots at after folder's latest shot at date, update it
			if folder.LatestShotAt == nil || ff.ShotAt.Compare(*folder.LatestShotAt) > 0 {
				folder.LatestShotAt = &ff.ShotAt
			}
			var count int64
			err = tx.Model(&FolderFile{}).Where("folder_id = ?", folder.ID).Count(&count).Error
			if err != nil {
				return err
			}
			folder.TotalFiles = count
			tx.Save(&folder)
		}
	}

	return nil
}

// Checks whether folderFile is currently being processed by a cloud upload job
func (ff *FolderFile) IsProcessingByCloudUpload() bool {
	for _, item := range []string{
		ff.HwOriginalFilePath,
		ff.GoogleFilePath,
		ff.GoogleOriginalFilePath,
		ff.GoogleThumbnailPath,
		ff.GoogleDriveFileID,
	} {
		if item == "processing" {
			return true
		}
	}
	return false
}

// Checks if folder file is an image file
func (ff *FolderFile) IsImage() bool {
	return utils.FileExtToFileType(filepath.Ext(ff.CompressedFilePath)) == "image"
}

// Checks if folder file is a video file
func (ff *FolderFile) IsVideo() bool {
	return utils.FileExtToFileType(filepath.Ext(ff.CompressedFilePath)) == "video"
}

// Get the thumbnail path
func (ff *FolderFile) GenerateThumbnailPath() string {
	return utils.ChangeFileExtension(ff.CompressedFilePath, "jpg") + ".thumbnail.jpg"
}

func (ff FolderFile) MarshalJSON() ([]byte, error) {
	// Define a temporary struct to hold the marshalled data
	type FolderFileMarshal struct {
		BaseDbModel
		PeerID                 int64          `gorm:"column:peer_id" json:"peerId"`
		OwnerID                int64          `gorm:"column:owner_id" json:"ownerId"`
		FolderID               int64          `gorm:"column:folder_id" json:"folderId"`
		Disk                   string         `gorm:"column:disk" json:"disk"`
		FileName               string         `gorm:"column:file_name" json:"fileName"`
		FileType               string         `gorm:"column:file_type" json:"fileType"`
		OriginalFilePath       string         `gorm:"column:original_file_path" json:"originalFilePath"`
		CompressedFilePath     string         `gorm:"column:file_path" json:"filePath"`
		ThumbnailPath          string         `gorm:"column:thumbnail_path" json:"thumbnailPath"`
		HwOriginalFilePath     string         `gorm:"column:hw_original_file_path" json:"hwOriginalFilePath"`
		GoogleOriginalFilePath string         `gorm:"column:google_original_file_Path" json:"googleOriginalFilePath"`
		GoogleFilePath         string         `gorm:"column:google_file_path" json:"googleFilePath"`
		GoogleThumbnailPath    string         `gorm:"column:google_thumbnail_path" json:"googleThumbnailPath"`
		GoogleDriveFileID      string         `gorm:"column:google_drive_file_id" json:"googleDriveFileId"`
		Remark                 string         `gorm:"column:remark" json:"remark"`
		Metadata               map[string]any `gorm:"column:metadata;serializer:json" json:"-"`
		HasExif                bool           `gorm:"column:has_exif; default:0" json:"hasExif"`
		IsPrivate              bool           `gorm:"column:is_private; default:0" json:"isPrivate"`
		Views                  int64          `gorm:"column:views; default:0" json:"views"`
		SyncedAt               *time.Time     `gorm:"column:synced_at; null" json:"syncedAt"`
		IsDownloading          bool           `gorm:"column:is_downloading; default:0" json:"isDownloading"`
		ShotAt                 time.Time      `gorm:"column:shot_at" json:"-"`
		ShotAtString           string         `gorm:"-" json:"shotAt"`
		IsPreprocessing        int            `gorm:"-" json:"isPreprocessing"`
		MetadataSimple         map[string]any `gorm:"-" json:"metadata"`
	}

	shotAtString := ff.ShotAt.UTC().Format(time.RFC3339)
	if len(shotAtString) > 0 {
		microsecond := strconv.Itoa(int(ff.ShotAt.UTC().UnixMicro()))
		microsecond = microsecond[len(microsecond)-6:]
		shotAtString = shotAtString[:len(shotAtString)-1] + "." + microsecond + shotAtString[len(shotAtString)-1:]
	}
	ff.ShotAtString = shotAtString

	isChina := os.Getenv("APP_CHINA") == "true"
	if isChina {
		if ff.OriginalFilePath == "" {
			ff.IsPreprocessing = 1
		}
	} else {
		if ff.HwOriginalFilePath == "" || ff.GoogleFilePath == "" {
			ff.IsPreprocessing = 1
		}
	}

	ff.MetadataSimple = make(map[string]any)
	ff.MetadataSimple["dimension"] = ff.Metadata["dimension"]
	ff.MetadataSimple["duration"] = ff.Metadata["duration"]
	ff.MetadataSimple["size"] = ff.Metadata["size"]

	return json.Marshal(FolderFileMarshal(ff))
}
