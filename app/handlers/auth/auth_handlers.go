package auth

import (
	"net/http"
	"os"

	"github.com/innovay-software/famapp-main/app/dto"
	apiErrors "github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/handlers"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"pageTitle":       "Login",
		"loginUrl":        "https://" + c.Request.Host + "/api/v2/oauth/login",
		"refreshTokenUrl": "https://" + c.Request.Host + "/api/v2/oauth/refresh-token",
	})
}

func MobileLoginHandler(
	c *gin.Context, mobile, password, deviceToken string,
) (
	dto.ApiResponse, error,
) {
	user, err := repositories.UserRepoIns.FindUserByField("mobile", mobile)
	if err != nil {
		return nil, apiErrors.ApiErrorCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// return nil, errors.ApiErrorCredentials
	}

	tokenString, err := services.GenerateJwtAccessToken(user.UUID.String())
	if err != nil {
		return nil, apiErrors.ApiErrorSystem
	}

	return generateUserLoggedInData(c, user, tokenString, deviceToken), nil
}

func AccessTokenLoginHandler(
	c *gin.Context, user *models.User, deviceToken string,
) (
	dto.ApiResponse, error,
) {
	user.DeviceToken = &deviceToken
	if err := repositories.SaveDbModel(user); err != nil {
		return nil, err
	}

	tokenString, err := services.GenerateJwtAccessToken(user.UUID.String())
	if err != nil {
		return nil, apiErrors.ApiErrorSystem
	}

	return generateUserLoggedInData(c, user, tokenString, deviceToken), nil
}

func generateUserLoggedInData(
	c *gin.Context, user *models.User, accessToken, deviceToken string,
) *dto.LoginResponse {

	refreshToken := ""
	if tokenString, err := services.GenerateJwtRefreshToken(user.UUID.String()); err == nil {
		refreshToken = tokenString
	}
	handlers.SetAccessAndRefreshTokens(accessToken, refreshToken, c)

	// Set DeviceToken
	if deviceToken != "" {
		user.DeviceToken = &deviceToken
	}

	// Update user to DB
	repositories.SaveDbModel(user)
	// Sync user with mongo DB
	repositories.MongoRepoIns.CreateUserInMongo(user)
	folders := repositories.UserRepoIns.FindFolders(user)
	for _, folder := range folders {
		folder.Invitees = repositories.FolderRepoIns.FindInvitees(folder)
	}

	// Generate LoginResponse instance
	res := dto.LoginResponse{
		ApiResponseBase: dto.ApiResponseBase{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		GoogleCloudStorageAccessId:        os.Getenv("GOOGLE_CLOUD_STORAGE_ACCESS_ID"),
		GoogleCloudStorageAccessSecretKey: os.Getenv("GOOGLE_CLOUD_STORAGE_ACCESS_SECRET_KEY"),
		GoogleCloudStorageDomain:          os.Getenv("GOOGLE_CLOUD_STORAGE_DOMAIN"),
		GoogleCloudStorageBucketName:      os.Getenv("GOOGLE_CLOUD_STORAGE_BUCKET_NAME"),
		HwObsAccessId:                     os.Getenv("HWY_ACCESS_ID"),
		HwObsAccessSecretKey:              os.Getenv("HWY_ACCESS_SECRET_KEY"),
		HwObsDomain:                       os.Getenv("HWY_ENDPOINT"),
		HwObsBucketName:                   os.Getenv("HWY_BUCKET_NAME"),
		User:                              user,
		Folders:                           folders,
	}

	if allMembers, err := repositories.UserRepoIns.FindMemberList(1000, 0); err == nil {
		res.Members = allMembers
	}

	// if folders, err := repositories.FolderRepoIns.GetUserFolders(user.ID); err == nil {
	// 	res.Folders = folders
	// }

	return &res
}
