package handlers

// import (
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/innovay-software/famapp-main/app/errors"
// 	"github.com/innovay-software/famapp-main/app/models"
// 	"github.com/innovay-software/famapp-main/app/repositories"
// 	"github.com/innovay-software/famapp-main/app/services"
// )

// // Public Api handler wrapper
// func ApiHandlerWrapper(c *gin.Context, handler ApiHandler) bool {
// 	data, err := handler(c)
// 	if err != nil {
// 		return apiRespFailError(err, data, c)
// 	}
// 	return apiRespSucess(data, c)
// }

// // Public Api File Handler wrapper
// func ApiFileHandlerWrapper(c *gin.Context, handler ApiFileHandler) bool {
// 	if err := handler(c); err != nil {
// 		return apiRespFailError(err, nil, c)
// 	}
// 	return true
// }

// // Private User Api handler wrapper
// func ApiUserHandlerWrapper(c *gin.Context, handler ApiUserHandler) bool {
// 	user, err := getAuthenticatedUser(c)
// 	if err != nil {
// 		return apiRespFailError(err, nil, c)
// 	}

// 	data, err := handler(c, user)
// 	if err != nil {
// 		return apiRespFailError(err, data, c)
// 	}
// 	return apiRespSucess(data, c)
// }

// // Private File Display handler wrapper
// func ApiUserFileHandlerWrapper(c *gin.Context, handler ApiUserFileHandler) bool {
// 	user, err := getAuthenticatedUser(c)
// 	if err != nil {
// 		return apiRespFailError(err, nil, c)
// 	}
// 	if err := handler(c, user); err != nil {
// 		return apiRespFailError(err, nil, c)
// 	}
// 	return true
// }

// // Private Admin Api handler wrapper
// func ApiAdminHandlerWrapper(c *gin.Context, handler ApiAdminHandler) bool {
// 	user, err := getAuthenticatedAdmin(c)
// 	if err != nil {
// 		return apiRespFailError(err, nil, c)
// 	}

// 	data, err := handler(c, user)
// 	if err != nil {
// 		return apiRespFailError(err, data, c)
// 	}
// 	return apiRespSucess(data, c)
// }

// // Authenticate Admin
// func getAuthenticatedAdmin(c *gin.Context) (*models.User, error) {
// 	user, err := getAuthenticatedUser(c)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !user.IsAdmin() {
// 		return nil, errors.ApiErrorPermissionDenied
// 	}

// 	return user, nil
// }

// // Authenticate User
// func getAuthenticatedUser(c *gin.Context) (*models.User, error) {
// 	// Extract authentication token from either the header or query param
// 	tokenString := getAuthenticationToken(c)
// 	tokenString = strings.Trim(tokenString, " ")
// 	if tokenString == "" {
// 		return nil, errors.ApiErrorToken
// 	}

// 	// Verify the integrity of JwtToken
// 	jwtToken, err := services.VerifyToken(tokenString)
// 	if err != nil {
// 		return nil, errors.ApiErrorToken
// 	}

// 	// Get JWT Sub
// 	sub, err := jwtToken.Claims.GetSubject()
// 	if err != nil || sub == "" {
// 		return nil, errors.ApiErrorToken
// 	}

// 	// Get user
// 	user, err := repositories.UserRepoIns.FindUserByField("uuid", sub)
// 	if err != nil {
// 		return nil, errors.ApiErrorCredentials
// 	}

// 	newAccessToken, newRefreshToken := "", ""

// 	// Check for expiration
// 	expireTime, err := jwtToken.Claims.GetExpirationTime()
// 	if err != nil || expireTime.Time.UTC().Compare(time.Now().UTC()) == -1 {
// 		// If can't get expire time, or expire time is before current time (expired),
// 		// Generate a new token using refresh token

// 		// Empty token string since it already expired
// 		oldRefreshToken := c.Request.Header["Refresh-Token"]
// 		refreshToken, err := verifyRefreshToken(user, oldRefreshToken[0])
// 		if err == nil && len(refreshToken) > 0 {
// 			newJwtToken, err := services.GenerateJwtToken(user.Mobile)
// 			if err == nil {
// 				newAccessToken = newJwtToken
// 				newRefreshToken = refreshToken
// 			}
// 		}
// 		if newAccessToken == "" {
// 			// if tokenString wasn't successfully regenerated, return error
// 			return nil, errors.ApiErrorToken
// 		}
// 	}

// 	setAccessAndRefreshTokens(newAccessToken, newRefreshToken, c)
// 	return user, nil
// }

// // Extract JWT token from either the header or query params
// func getAuthenticationToken(c *gin.Context) string {
// 	tokenHeader := c.Request.Header["Authorization"]
// 	if len(tokenHeader) == 1 {
// 		tokenString := strings.Replace(tokenHeader[0], "Bearer ", "", 1)
// 		return tokenString
// 	}
// 	token, exists := c.GetQuery("token")
// 	if exists {
// 		return token
// 	}
// 	return ""
// }

// // Refresh authentication token based on refresh token
// func verifyRefreshToken(user *models.User, refreshToken string) (string, error) {
// 	// Todo: Check for black listed refreshToken
// 	// End of Todo
// 	return refreshToken, nil
// }
