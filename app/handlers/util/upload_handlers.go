package util

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
	"github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/services"
	"github.com/innovay-software/famapp-main/app/utils"
)

func Base64ChunkUploadFileHandler(
	c *gin.Context, user *models.User,
	base64Content, fileName, chunkedFileName string, hasMore bool,
) (
	dto.ApiResponse, error,
) {

	byteList, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return nil, errors.ApiErrorParamInvalid
	}

	ext := filepath.Ext(fileName)
	if chunkedFileName == "" {
		chunkedFileName = fileName + "-" + utils.GenerateRandomString(10, true, false, false) + ext
	}

	chunkFileAbsPath := utils.GetStorageAbsPath("chunk-upload", chunkedFileName)

	f, err := os.OpenFile(chunkFileAbsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.ApiErrorSystem
	}
	defer f.Close()

	if _, err := f.Write(byteList); err != nil {
		return nil, errors.ApiErrorSystem
	}

	if hasMore {
		return &dto.Base64ChunkUploadFileResponse{
			RemoteFileId:    0,
			Uploaded:        true,
			ChunkedFileName: chunkedFileName,
			Document:        nil,
			HasMore:         true,
		}, nil
	}

	metadataMap := services.ExtractFileMetadata(chunkFileAbsPath)
	targetDatetime := (*metadataMap)["shot_at_date_time"].(time.Time)

	relativeDirPath := strings.Join([]string{
		strconv.FormatInt(user.ID, 10),
		strings.ReplaceAll(targetDatetime.Format(time.DateOnly), "-", "/"),
	}, "/")
	fileNamePostFix := strconv.Itoa(rand.Intn(8888) + 1001)
	// fmt.Println("fileNamePostFix = ", fileNamePostFix)
	fileName = strings.ReplaceAll(targetDatetime.Format(time.DateOnly), "-", "_") + "_" + fileNamePostFix + ext
	relativeFilePath := relativeDirPath + "/" + fileName

	uploadDisk := "user-upload"
	absDirPath := utils.GetStorageAbsPath(uploadDisk, relativeDirPath)
	if err := os.MkdirAll(absDirPath, 0755); err != nil {
		return nil, fmt.Errorf("unable to create dir")
	}

	absFilePath := absDirPath + "/" + fileName
	os.Rename(chunkFileAbsPath, absFilePath)

	fmt.Println("absFilePath", absFilePath)
	uploadIns := models.Upload{
		UserID:   user.ID,
		Disk:     uploadDisk,
		FileName: fileName,
		FileType: utils.FileExtToFileType(ext),
		FilePath: relativeFilePath,
		ShotAt:   targetDatetime.UTC(),
	}

	repositories.SaveDbModel(&uploadIns)

	return &dto.Base64ChunkUploadFileResponse{
		RemoteFileId:    uploadIns.ID,
		Uploaded:        true,
		ChunkedFileName: chunkedFileName,
		Document:        &uploadIns,
		HasMore:         false,
	}, nil
}