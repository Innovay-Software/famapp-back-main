package utils

import (
	"crypto/md5"
	"encoding/hex"
	"image"
	"io"
	"os"
	"slices"
	"strings"
)

// Checks if path exists
func PathExists(absolutePath string) bool {
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Make a copy of the file
func DuplicateFile(srcPath, dstPath string) error {
	in, err := os.Open(srcPath)
	if err != nil {
		LogError("Unable to open file:", srcPath, err)
		return err
	}
	defer in.Close()

	out, err := os.Create(dstPath)
	if err != nil {
		LogError("Unable to create file:", dstPath, err)
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		LogError("Unable to copy file:", out, in, err)
		return err
	}
	outErr := out.Sync()
	if outErr != nil {
		LogError("Unable to out.Synce():", outErr)
		return outErr
	}
	return nil
}

// Calculate MD5 for target file
func GenerateFileMd5(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := md5.New()
	_, err = io.Copy(w, f)
	if err != nil {
		return "", err
	}

	rawHash := w.Sum(nil)
	return hex.EncodeToString(rawHash), nil
}

// Get file type based on file extension
func FileExtToFileType(ext string) string {
	if ext[0] == '.' {
		ext = ext[1:]
	}
	ext = strings.ToLower(ext)

	typeToExtMap := map[string]string{
		"image": "jpg,jpeg,png,gif,ico,svg,tiff",
		"video": "mp4,mov",
		"pdf":   "pdf",
		"doc":   "doc,docx",
		"excel": "xls,xlsx,xlt,xltx,xlsm",
	}

	for k, v := range typeToExtMap {
		if slices.Contains(strings.Split(v, ","), ext) {
			return k
		}
	}
	return "others"
}

// Get file size
func GetFileSize(filepath string) int64 {
	fi, err := os.Stat(filepath)
	if err == nil {
		LogError("utils.GetFileSize: Unabled to stat file:" + filepath)
		return fi.Size()
	}
	return 0
}

// Get image dimension
func GetImageDimeision(imageAbsolutePath string) (int, int) {
	file, err := os.Open(imageAbsolutePath)
	if err != nil {
		return 0, 0
	}

	imageConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0
	}

	return imageConfig.Width, imageConfig.Height
}

// Change the file extension
func ChangeFileExtension(filePath, newExtension string) string {
	n := len(filePath)
	if newExtension[0] != '.' {
		newExtension = "." + newExtension
	}
	for i := n - 1; i >= 0; i-- {
		if filePath[i] == '.' {
			return filePath[:i] + newExtension
		}
	}
	return filePath
}

// Delete file from disk
func DeleteFile(absPath string) error {
	if PathExists(absPath) {
		if err := os.Remove(absPath); err != nil {
			LogError("Unable to remove file:", absPath)
			return err
		}
	}
	return nil
}
