package utils

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/innovay-software/famapp-main/config"
	"golang.org/x/image/draw"
)

// Compressed to jpg format with max(width, height) = targetSize
// Returns the result image path and error
func CompressImageToJpgWithMaxSize(
	imageAbsolutePath string,
	targetSize int,
) (
	string, error,
) {
	srcWidth, srcHeight := GetImageDimeision(imageAbsolutePath)
	if srcWidth == 0 || srcHeight == 0 {
		return "", errors.New("Cannot get image dimension: " + imageAbsolutePath)
	}

	// Get target width and height
	dstWidth, dstHeight := srcWidth, srcHeight
	if srcWidth > targetSize || srcHeight > targetSize {
		if srcWidth > srcHeight {
			// horizontal
			dstWidth, dstHeight = targetSize, targetSize*srcHeight/srcWidth
		} else {
			// vertical
			dstWidth, dstHeight = targetSize*srcWidth/srcHeight, targetSize
		}
	}

	return compressImageToJpgWithTargetDimension(imageAbsolutePath, dstWidth, dstHeight)
}

// Compressed to jpg format with min(width, height) = targetSize
// Returns the result image path and error
func CompressImageToJpgWithMinSize(
	imageAbsolutePath string,
	targetSize int,
) (
	string, error,
) {
	srcWidth, srcHeight := GetImageDimeision(imageAbsolutePath)
	if srcWidth == 0 || srcHeight == 0 {
		return "", errors.New("Cannot get image dimension: " + imageAbsolutePath)
	}

	// Get target width and height
	dstWidth, dstHeight := srcWidth, srcHeight
	if srcWidth > targetSize || srcHeight > targetSize {
		if srcWidth < srcHeight {
			// horizontal
			dstWidth, dstHeight = targetSize, targetSize*srcHeight/srcWidth
		} else {
			// vertical
			dstWidth, dstHeight = targetSize*srcWidth/srcHeight, targetSize
		}
	}

	return compressImageToJpgWithTargetDimension(imageAbsolutePath, dstWidth, dstHeight)
}

func compressImageToJpgWithTargetDimension(
	imageAbsolutePath string, dstWidth, dstHeight int,
) (
	string, error,
) {
	// Use a temp location to store intermediate file
	tempImagePath := imageAbsolutePath + ".temp.jpg"
	resultImagePath := ChangeFileExtension(imageAbsolutePath, "jpg")

	// check if image file exist
	if !PathExists(imageAbsolutePath) {
		LogError("File not found:", imageAbsolutePath)
		return "", fmt.Errorf("file does not exist")
	}

	// Read source image
	srcImage, _, _, err := readImage(imageAbsolutePath)
	if err != nil {
		return "", err
	}

	// Create destination image
	dstImage := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))

	// Draw destination image
	draw.NearestNeighbor.Scale(dstImage, dstImage.Rect, srcImage, srcImage.Bounds(), draw.Over, nil)

	// Write destination image
	if err := writeImageAsJpg(dstImage, tempImagePath); err != nil {
		Log("writeImageAsJpg failed with error:", err)
		return "", err
	}

	os.Remove(resultImagePath)
	os.Rename(tempImagePath, resultImagePath)

	return resultImagePath, nil
}

// // Resize the dimension proportionally with width and height not greater than maxSize
// // If width and height is less than maxSize, they're returned
// // If width and height is greater than maxSize, they're scaled down
// func ScaleDownWidthHeightProportionally(width, height, maxSize int) (int, int) {
// 	// If both width and height are less than maxSize, no need to resize, return original size
// 	if width <= maxSize && height <= maxSize {
// 		return width, height
// 	}

// 	if width >= height {
// 		// dimension is horizontal
// 		return maxSize, maxSize * height / width
// 	}

// 	// dimension is vertical
// 	return maxSize * width / height, maxSize
// }

// // Resize the dimension proportionally with width and height at least targetSize
// func GetPreferredImageDimensionWithTargetSize(width, height, targetSize int) (int, int) {
// 	if width == height {
// 		return targetSize, targetSize
// 	}
// 	if width < height {
// 		// image is horizontal, set width to maxSize and height shrinks proportionally
// 		// Multiply first to increase acuracy
// 		// targetWidth, targetHeight = targetSize, (targetSize * height / width)
// 		return targetSize, targetSize * height / width
// 	}
// 	// iamge is vertical, set height to maxSize and width shrinks proportionally
// 	// Multiply first to increase acuracy
// 	// targetWidth, targetHeight = (targetSize * width / height), targetSize
// 	return targetSize * width / height, targetSize
// }

func GenerateThumbnailJpg(
	srcImageAbsPath, thumbnailAbsPath string,
) (
	string, error,
) {
	thumbnailAbsPath = ChangeFileExtension(thumbnailAbsPath, "jpg")
	if !PathExists(srcImageAbsPath) {
		LogError("Missing source image file:", srcImageAbsPath)
		return "", errors.New("Missing source image file: " + srcImageAbsPath)
	}
	if PathExists(thumbnailAbsPath) {
		os.Remove(thumbnailAbsPath)
	}

	// Make a copy of source image to thumbnail path
	if err := DuplicateFile(srcImageAbsPath, thumbnailAbsPath); err != nil {
		return "", err
	}
	srcImageAbsPath = ""

	// Compress the thumbnail file to make it jpg with target size
	thumbnailAbsPathTemp, err := CompressImageToJpgWithMinSize(
		thumbnailAbsPath, config.FolderFileThumbnailSize,
	)
	if err != nil {
		return "", err
	}
	thumbnailAbsPath = thumbnailAbsPathTemp

	// Crop the thumbnail file and save it
	// Read source image
	srcImage, srcWidth, srcHeight, err := readImage(thumbnailAbsPath)
	if err != nil {
		return "", err
	}
	if srcWidth == srcHeight {
		return thumbnailAbsPath, nil
	}

	if srcWidth > srcHeight {
		// Horizontal Image
		paddingLeft := (srcWidth - srcHeight) / 2
		srcImage2, err := cropImage(srcImage, image.Rect(paddingLeft, 0, paddingLeft+srcHeight, srcHeight))
		if err != nil {
			return "", err
		}
		srcImage = srcImage2
	} else if srcHeight > srcWidth {
		// Vertical Image
		paddingTop := (srcHeight - srcWidth) / 2
		srcImage2, err := cropImage(srcImage, image.Rect(0, paddingTop, srcWidth, paddingTop+srcWidth))
		if err != nil {
			return "", err
		}
		srcImage = srcImage2
	}

	// Write destination image
	if err := writeImageAsJpg(srcImage, thumbnailAbsPath); err != nil {
		return "", err
	}
	return thumbnailAbsPath, nil
}

// Reads image from disk
func readImage(fileAbsPath string) (image.Image, int, int, error) {
	f, err := os.Open(fileAbsPath)
	if err != nil {
		LogError("Cannot open file:", fileAbsPath)
		return nil, 0, 0, err
	}
	defer f.Close()

	// Read source image file
	img, _, err := image.Decode(f)
	if err != nil {
		LogError("Cannot decode file to image")
		return nil, 0, 0, err
	}

	return img, img.Bounds().Max.X, img.Bounds().Max.Y, nil
}

// cropImage takes an image and crops it to the specified rectangle.
func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

// writeImage writes an Image back to the disk.
func writeImageAsJpg(img image.Image, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, &jpeg.Options{Quality: 80})
}
