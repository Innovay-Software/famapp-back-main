package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/redis/go-redis/v9"
)

var keyPrefix = ""

type messageQueueRepo struct {
	redisCon *redis.Client
}

func (mr *messageQueueRepo) getKeyPrefix() string {
	if keyPrefix == "" {
		keyPrefix = os.Getenv("REDIS_KEY_PREFIX")
	}
	return keyPrefix
}

func (mr *messageQueueRepo) userSyncQueueKey() string {
	return mr.getKeyPrefix() + ":UserInfoSync"
}

func (mr *messageQueueRepo) folderFileProcessingQueueKey() string {
	return mr.getKeyPrefix() + ":FolderFileProcessing"
}

func (mr *messageQueueRepo) failedFolderFileUploadQueueKey() string {
	return mr.getKeyPrefix() + ":FailedFolderFileUpload"
}

func (mr *messageQueueRepo) cloudUploadLastRunTimeKey() string {
	return mr.getKeyPrefix() + ":CloudUploadLastRunTime"
}

func (mr *messageQueueRepo) sendUserInfoToUserSyncQueue(user *models.User) error {
	if user == nil {
		return errors.New("missing user")
	}
	userData := map[string]string{
		"uuid":      user.UUID.String(),
		"family_id": fmt.Sprintf("%v", user.FamilyID),
		"name":      user.Name,
		"email":     user.Email,
		"mobile":    user.Mobile,
		"avatar":    user.Avatar,
	}
	userDataJsonString, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	if err := mr.redisCon.RPush(context.Background(), mr.userSyncQueueKey(), userDataJsonString).Err(); err != nil {
		utils.LogError(err)
		return err
	}

	return nil
}

func (mr *messageQueueRepo) sendFolderFileToFolderFileProcessingQueue(folderFile *models.FolderFile) error {
	if folderFile == nil {
		return errors.New("missing folderFile")
	}
	if err := mr.redisCon.RPush(context.Background(), mr.folderFileProcessingQueueKey(),  folderFile.ID).Err(); err != nil {
		utils.LogError(err)
		return err
	}

	return nil
}

func (mr *messageQueueRepo) LpopUserInfoFromUserSyncQueue() (
	string, error,
) {
	content, err := mr.redisCon.LPop(context.Background(), mr.userSyncQueueKey()).Result()
	if err != nil {
		utils.LogError(err)
		return "", err
	}

	return content, nil
}

func (mr *messageQueueRepo) LpopFolderFileIdFromFolderFileProcessingQueue() (
	string, error,
) {
	content, err := mr.redisCon.LPop(context.Background(), mr.folderFileProcessingQueueKey()).Result()
	if err != nil {
		utils.LogError(err)
		return "", err
	}

	return content, nil
}

func (mr *messageQueueRepo) RpushFailedFolderFileIdsToFailedFolderFileUploadQueue(
	folderFileIdStringList []string,
) {
	for _, item := range folderFileIdStringList {
		err := mr.redisCon.RPush(context.Background(), mr.failedFolderFileUploadQueueKey(), item).Err()
		if err != nil {
			utils.LogError(err)
		}
	}
}

func (mr *messageQueueRepo) SetLastRunTime() {
	err := mr.redisCon.Set(context.Background(), mr.cloudUploadLastRunTimeKey(), time.Now().Unix(), 0).Err()
	if err != nil {
		utils.LogError(err)
	}
}
