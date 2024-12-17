package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/innovay-software/famapp-main/app/utils"
)

func ExtractFileMetadata(fileAbsPath string) *map[string]any {
	// Basic metadata fields
	metadata := map[string]any{
		"dimension": "1920x1080",
		"size":      utils.GetFileSize(fileAbsPath),
		"exif":      nil,
	}

	// If is China server, do not processing images and videos (wait for remote server updates)
	isChinaServer := os.Getenv("APP_CHINA")
	if isChinaServer == "true" {
		return &metadata
	}

	ext := filepath.Ext(fileAbsPath)
	fileType := utils.FileExtToFileType(ext)
	if fileType == "image" {
		if exifDataMap, err := utils.ExtractImageExif(fileAbsPath); err == nil {
			metadata["exif"] = *exifDataMap
		}
	} else if fileType == "video" {
		if ffprobeDataMap, err := utils.ExtractVideoFfprobeData(fileAbsPath); err == nil {
			metadata["exif"] = *ffprobeDataMap
		}
	}

	shotAtDate := time.Now()
	if exifDataMap, exists := metadata["exif"]; exists {
		exifDataMapPointer, ok := exifDataMap.(map[string]any)
		if !ok {
			utils.LogError("Unabled to cast exifDataMap = map[string]any type")
		}
		extractedDateTime, err := extractDateTimeFromExifMap(&exifDataMapPointer)
		if err == nil {
			shotAtDate = extractedDateTime
		} else {
			utils.LogWarning("Warning:", err)
		}
	}

	metadata["shot_at_date_time"] = shotAtDate
	return &metadata
}

func extractDateTimeFromExifMap(exifMap *map[string]any) (time.Time, error) {
	// Image Exif Format:
	hourOffset := ""
	if v, exists := (*exifMap)["Offset Time For DateTime"]; exists {
		if v2, ok := v.(string); ok {
			hourOffset = strings.Replace(v2, ":", "", -1)
		}
	}
	if v, exists := (*exifMap)["Date and Time"]; exists {
		if v2, ok := v.(string); ok {
			dateString := strings.Split(v2, " ")[0]
			timeString := strings.Split(v2, " ")[1]
			year := dateString[0:4]
			month := dateString[5:7]
			day := dateString[8:10]
			dateString = strings.Join(([]string{year, month, day}), "-")
			dateTimeParseString := dateString + " " + timeString
			dateTimeParseLayout := "2006-01-02 15:04:05"

			if len(hourOffset) > 0 {
				dateTimeParseString += " " + hourOffset
				dateTimeParseLayout += " " + "-0700"
			}

			// fmt.Println("dateTimeParseString=", dateTimeParseString)
			return time.Parse(dateTimeParseLayout, dateTimeParseString)
		}
	}

	// Video FFProbe Format:
	if creationDateTime, exists := (*exifMap)["format.tags.com_apple_quicktime_creationdate"]; exists {
		utils.Log("Found createimt time: ", creationDateTime)
		dateTimeParseLayout := "2006-01-02T15:04:05-0700"
		return time.Parse(dateTimeParseLayout, creationDateTime.(string))
	}

	return time.Now(), fmt.Errorf("unable to find DateTime in Exif map")
}
