package utils

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Get video duration in seconds
func ExtractVideoDuration(
	videoAbsPath string,
) (
	int, error,
) {
	byteSlice, err := exec.Command(
		"ffprobe",
		"-i",
		videoAbsPath,
		"-v",
		"quiet",
		"-show_entries",
		"format=duration",
		"-hide_banner",
		"-of",
		"default=noprint_wrappers=1:nokey=1",
	).Output()
	if err != nil {
		return 0, err
	}
	outputString := strings.Replace(string(byteSlice), "\n", "", -1)
	duration, err := strconv.ParseFloat(outputString, 64)
	if err != nil {
		return 0, err
	}
	return int(math.Ceil(duration)), nil
}

// Extract the first frame from video and save to the path specified
// If the specified doesn't have a .jpg extension, it will be changed to .jpg
func ExtractVideoFirstFrameAsJpg(
	videoAbsPath string, frameSavePath string,
) (
	string, error,
) {
	// Extract first frame and save as jpg: $ ffmpeg -i input.mp4 -frames:v 1 output.jpg

	// Change to jpg extension
	frameSavePath = ChangeFileExtension(frameSavePath, "jpg")

	// Delete frame jpg file if it exists
	if err := DeleteFile(frameSavePath); err != nil {
		return "", err
	}

	// Generate frame jpg
	argumentList := []string{
		"-i", videoAbsPath, "-frames:v", "1", frameSavePath,
	}
	cmd := exec.Command("ffmpeg", argumentList...)
	out, err := cmd.CombinedOutput()
	if err != nil || !PathExists(frameSavePath) {
		LogError("Unabled to extract frame to:", frameSavePath, err)
		LogError("cmd result:")
		LogError(string(out))
		return "", err
	}
	return frameSavePath, nil
}

// Extract video FFProbe data
func ExtractVideoFfprobeData(
	fileAbsPath string,
) (
	*map[string]any, error,
) {
	if !PathExists(fileAbsPath) {
		LogError("Unable to get file stats:", fileAbsPath)
		return nil, errors.New("Missing file: " + fileAbsPath)
	}

	metadataMap := make(map[string]any)
	byteSlice, err := exec.Command(
		"ffprobe",
		"-loglevel",
		"error",
		"-show_entries",
		"stream_tags:format_tags",
		"-of",
		"flat",
		fileAbsPath,
	).Output()
	if err != nil {
		LogError("ffprobe command line error:", err)
		return nil, err
	}

	exifString := string(byteSlice)
	for _, line := range strings.Split(exifString, "\n") {
		lineParts := strings.Split(line, "=")
		if len(lineParts) == 2 {
			vString := lineParts[1]
			metadataMap[lineParts[0]] = vString[1:len(vString)-1]
		}
	}

	return &metadataMap, nil
}

// Compressed video to mp4 fullHD dimension
func CompressVideoToMp4FullHD(
	videoAbsPath string,
) (
	string, error,
) {
	if !PathExists(videoAbsPath) {
		LogError("File not exist:", videoAbsPath)
		return "", errors.New("Missing file: " + videoAbsPath)
	}

	tempFileAbsPath := videoAbsPath + ".temp.mp4"
	firstFramePath := videoAbsPath + ".frame.jpg"
	firstFramePath, _ = ExtractVideoFirstFrameAsJpg(videoAbsPath, firstFramePath)
	defer os.Remove(firstFramePath)

	width, height := GetImageDimeision(firstFramePath)
	targetWidth, targetHeight := width, height
	if width > height {
		targetWidth, targetHeight = 1280, 1280*height/width
	} else {
		targetWidth, targetHeight = 1280*width/height, 1280
	}

	// Command to resize: $ffmpeg -i input.avi -vf scale=320:240 output.avi
	newScale := fmt.Sprintf("%d:%d", targetWidth, targetHeight)
	argumentList := []string{
		"-i", videoAbsPath, "-vf", "scale=" + newScale, tempFileAbsPath,
	}
	cmd := exec.Command("ffmpeg", argumentList...)
	out, err := cmd.CombinedOutput()
	if err != nil || !PathExists(tempFileAbsPath) {
		LogError("Ffmpeg command failed:", err)
		LogError(string(out))
		return "", errors.New("unable to extract video frist frame")
	}

	if err := DeleteFile(videoAbsPath); err != nil {
		LogError("Unable to delete file", videoAbsPath, err)
		return "", err
	}

	videoAbsPath = ChangeFileExtension(videoAbsPath, "mp4")
	if err := os.Rename(tempFileAbsPath, videoAbsPath); err != nil {
		LogError("Unable to rename file from", tempFileAbsPath, "to", videoAbsPath)
		return "", err
	}
	return videoAbsPath, nil
}
