package util

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
	"github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/utils"
)

// Get a single config value. If it's not found, return ""
func GetConfig(c *gin.Context, user *models.User, configKey string) (
	dto.ApiResponse, error,
) {
	config := repositories.ConfigRepoIns.GetConfig(configKey)
	if config == nil {
		return &dto.GetConfigResponse{ConfigValue: ""}, nil
	}
	return &dto.GetConfigResponse{ConfigValue: config.ConfigValue}, nil
}

func Ping(c *gin.Context) (
	dto.ApiResponse, error,
) {
	return &dto.PingResponse{Ping: "pong"}, nil
}

func CheckForMobileUpdate(c *gin.Context, clientOs, currentVersion string) (
	dto.ApiResponse, error,
) {
	res := dto.CheckForUpdateResponse{
		HasUpdate:   false,
		ForceUpdate: false,
		Title:       "",
		Content:     map[string]any{},
		Url:         "",
	}

	if latestVersion, err := repositories.UtilsRepoIns.GetLatestAppVersion(clientOs); err == nil {
		res.Version = latestVersion.Version
		res.Title = latestVersion.Title
		res.Content = latestVersion.Content
		res.Url = latestVersion.Link

		// cOs := req.Os
		cVersion := currentVersion
		lVersion := latestVersion.Version
		cVersionParts := strings.Split(cVersion, ".")
		lVersionParts := strings.Split(lVersion, ".")
		if len(cVersionParts) < 1 || len(lVersionParts) < 1 {
			return nil, errors.ApiErrorSystem
		}

		// First number is the major version
		// if current major version is less than latest major version
		// then a force update is required
		cMajor, cMajErr := strconv.Atoi(cVersionParts[0])
		lMajor, lMajErr := strconv.Atoi(lVersionParts[0])
		if cMajErr != nil || lMajErr != nil {
			return nil, errors.ApiErrorSystem
		}

		if cMajor < lMajor {
			// If current major is behind latest major, a force update is required
			res.ForceUpdate = true
			res.HasUpdate = true
			return &res, nil
		} else if cMajor > lMajor {
			// If current major is ahead of latest major, no update is required
			res.ForceUpdate = false
			res.HasUpdate = false
			return &res, nil
		}

		// When cMajor == lMajor, compare minor versions:

		// Compare the rest of the version numbers, if any one from current version
		// is less than the corresponding number from latest version,
		// then an update is available but not required
		for i := 1; i < len(cVersionParts) && i < len(lVersionParts); i++ {
			cMinor, cMinErr := strconv.Atoi(cVersionParts[i])
			lMinor, lMinErr := strconv.Atoi(lVersionParts[i])
			if cMinErr != nil || lMinErr != nil {
				utils.LogError("Unable to compare minor versions:", currentVersion, latestVersion.Version)
				break
			}
			if cMinor < lMinor {
				res.HasUpdate = true
				res.ForceUpdate = false
				break
			}
		}
	}
	return &res, nil
}

func UserAvatar(c *gin.Context, userId int64) error {
	user, err := repositories.UserRepoIns.FindUserByField("id", strconv.Itoa(int(userId)))
	if err == nil && user != nil && user.Avatar != "" {
		c.Redirect(302, user.Avatar)
		return nil
	}
	c.Redirect(302, utils.GetUrlPath("static", "default-avatar.png"))
	return nil
}
