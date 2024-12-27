package utils

import (
	"os/exec"
	"strconv"
	"strings"
)

var (
	orientationStringToIntegerMap = map[string]int{
		"top-left":     1,
		"top-right":    2,
		"bottom-right": 3,
		"bottom-left":  4,
		"left-top":     5,
		"right-top":    6,
		"right-bottom": 7,
		"left-bottom":  8,
	}
)

// Extract EXIF data from image
func ExtractImageExif(fileAbsPath string) (*map[string]any, error) {
	exifMap := make(map[string]any)

	exifByteSlice, err := exec.Command("exif", "-m", fileAbsPath).Output()
	if err != nil {
		return &exifMap, err
	}

	exifString := string(exifByteSlice)
	for _, line := range strings.Split(exifString, "\n") {
		lineParts := strings.Split(line, "\t")
		if len(lineParts) == 2 {
			exifMap[lineParts[0]] = lineParts[1]
		}
	}
	return &exifMap, nil
}

func SetImageExif(fileAbsPath string, orientation string) error {
	Log("SetImageExif:", fileAbsPath, orientation, "|")
	if orientation != "" {
		// Create exif data
		_, err1 := exec.Command(
			"exif",
			"--output="+fileAbsPath,
			"--create-exif",
			fileAbsPath,
		).Output()
		if err1 != nil {
			return err1
		}

		orientationVal := 0
		if tVal, exists := orientationStringToIntegerMap[strings.ToLower(orientation)]; exists {
			orientationVal = tVal
		}
		_, err2 := exec.Command(
			"exif",
			"--output="+fileAbsPath,
			"--ifd=0",
			"--tag=Orientation",
			"--set-value="+strconv.Itoa(orientationVal),
			"--no-fixup",
			fileAbsPath,
		).Output()
		if err2 != nil {
			return err2
		}
	}
	return nil
}

// exif --output=2024_02_03_9893.jpg.thumbnail0.jpg --ifd=0 --tag=Orientation --set-value=0 --no-fixup 2024_02_03_9893.jpg.thumbnail.jpg
// exif --output=2024_02_03_9893.jpg.thumbnail1.jpg --ifd=0 --tag=Orientation --set-value=1 --no-fixup 2024_02_03_9893.jpg.thumbnail.jpg
// exif --output=2024_02_03_9893.jpg.thumbnail2.jpg --ifd=0 --tag=Orientation --set-value=2 --no-fixup 2024_02_03_9893.jpg.thumbnail.jpg
// exif --output=2024_02_03_9893.jpg.thumbnail3.jpg --ifd=0 --tag=Orientation --set-value=3 --no-fixup 2024_02_03_9893.jpg.thumbnail.jpg
// exif --output=2024_02_03_9893.jpg.thumbnail4.jpg --ifd=0 --tag=Orientation --set-value=4 --no-fixup 2024_02_03_9893.jpg.thumbnail.jpg
// exif --output=2024_02_03_9893.jpg.thumbnail5.jpg --ifd=0 --tag=Orientation --set-value=5 --no-fixup 2024_02_03_9893.jpg.thumbnail.jpg
// exif --output=2024_02_03_9893.jpg.thumbnail6.jpg --ifd=0 --tag=Orientation --set-value=6 --no-fixup 2024_02_03_9893.jpg.thumbnail.jpg
// exif --output=2024_02_03_9893.jpg.thumbnail7.jpg --ifd=0 --tag=Orientation --set-value=7 --no-fixup 2024_02_03_9893.jpg.thumbnail.jpg
// exif --output=2024_02_03_9893.jpg.thumbnail8.jpg --ifd=0 --tag=Orientation --set-value=8 --no-fixup 2024_02_03_9893.jpg.thumbnail.jpg
// exif --output=2024_02_03_9893.jpg.thumbnail9.jpg --ifd=0 --tag=Orientation --set-value=9 --no-fixup 2024_02_03_9893.jpg.thumbnail.jpg

// exif --output=2024_02_03_9893.jpg.thumbnail.jpg --create-exif 2024_02_03_9893.jpg.thumbnail.jpg
