package middlewares

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	apiErrors "github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/services"
	"github.com/innovay-software/famapp-main/config"
)

// JWT bearerAuth token validation
// Checks if token is required, and refreshes token if token is about to expire
func JwtAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtAuthenticationErrorAction := func(enforceJwt bool) {
			if enforceJwt {
				panic(apiErrors.ApiErrorToken)
			} else {
				c.Next()
			}
		}

		// This header is set by the oapi-codegen's auto generate auth method
		// If a path requires authentication, it sets this header in the request
		// so this jwt middleware will be notified to do the authentication
		requireJwtHeader := c.Request.Header[config.RequireJwtHeaderKey]
		enforcingJwtAuthentication := true
		if len(requireJwtHeader) == 0 || requireJwtHeader[0] == "" {
			// RequireJwtHeaderKey is not found
			enforcingJwtAuthentication = false
		}

		// Get access token and refresh token
		accessTokenString := getAccessToken(c)
		accessToken, accessTokenValid, accessTokenExpired := inspectJwtToken(accessTokenString)
		if !accessTokenValid && !accessTokenExpired {
			jwtAuthenticationErrorAction(enforcingJwtAuthentication)
			return
		}
		refreshTokenString := getRefreshToken(c)
		refreshToken, refreshTokenValid, refreshTokenExpired := inspectJwtToken(refreshTokenString)

		newAccessTokenString, newRefreshTokenString := accessTokenString, refreshTokenString

		// Get JWT Sub
		sub, err := accessToken.Claims.GetSubject()
		if err != nil || sub == "" {
			jwtAuthenticationErrorAction(enforcingJwtAuthentication)
			return
		}

		// Get user
		user, err := repositories.UserRepoIns.FindUserByField("uuid", sub)
		if err != nil {
			jwtAuthenticationErrorAction(enforcingJwtAuthentication)
			return
		}
		c.Set(config.AuthenticatedUserHeaderKey, user)

		// If refresh token is provided and is valid, check if it needs to be refreshed
		if refreshTokenValid {
			hoursUntilExpiry := getHoursUntilExpiry(refreshToken)
			if hoursUntilExpiry < 10*24 {
				tokenString, err := services.GenerateJwtRefreshToken(user.UUID.String())
				if err == nil {
					newRefreshTokenString = tokenString
					refreshTokenValid = true
					refreshTokenExpired = false
				}
			}
		}

		// Handle expired access token
		if accessTokenExpired {
			if refreshTokenExpired {
				jwtAuthenticationErrorAction(enforcingJwtAuthentication)
				return
			} else {
				// if refresh token is not expired, generate a new access token
				tokenString, err := services.GenerateJwtAccessToken(user.UUID.String())
				if err != nil {
					jwtAuthenticationErrorAction(enforcingJwtAuthentication)
					return
				}
				newAccessTokenString = tokenString
				accessTokenValid = true
				accessTokenExpired = false
			}
		}

		// If access token was not updated, empty it
		if newAccessTokenString == accessTokenString {
			newAccessTokenString = ""
		}

		// If refresh token was not updated, empty it
		if newRefreshTokenString == refreshTokenString {
			newRefreshTokenString = ""
		}

		setAccessAndRefreshTokens(c, newAccessTokenString, newRefreshTokenString)
		c.Next()
	}
}

// Inspect a JWT token, returns (token, valid, expired)
func inspectJwtToken(tokenString string) (*jwt.Token, bool, bool) {
	if tokenString == "" {
		return nil, false, false
	}

	isTokenExpired := false
	jwtToken, err := services.VerifyToken(tokenString)
	if err != nil {
		if apiError, ok := err.(apiErrors.ApiError); ok {
			if apiError.Code == apiErrors.ApiErrorTokenExpired.Code {
				isTokenExpired = true
			}
		}
	}
	return jwtToken, err == nil, isTokenExpired
}

// Get number of hours until expiry. Negative values indicate already expired
func getHoursUntilExpiry(token *jwt.Token) int {
	expiryDate, err := token.Claims.GetExpirationTime()
	if err != nil {
		return -1
	}

	// duration := expiryDate.Sub(time.Now())
	duration := time.Until(expiryDate.Time)
	if duration.Minutes() < 0 {
		return -1
	}
	return int(duration.Minutes() / 60)
}

// Extract JWT authorization bearer token from either the header or query params
func getAccessToken(c *gin.Context) string {
	tokenHeader := c.Request.Header["Authorization"]
	if len(tokenHeader) == 1 {
		tokenString := strings.Replace(tokenHeader[0], "Bearer ", "", 1)
		return strings.Trim(tokenString, " ")
	}
	token, exists := c.GetQuery("token")
	if exists {
		return strings.Trim(token, " ")
	}
	return ""
}

// Extract JWT refresh token from the header
func getRefreshToken(c *gin.Context) string {
	tokenHeader := c.Request.Header["Refresh-Token"]
	if len(tokenHeader) == 1 {
		return tokenHeader[0]
	}
	return ""
}

// Set access token and refresh token to header
func setAccessAndRefreshTokens(c *gin.Context, accessToken, refreshToken string) {
	if accessToken != "" {
		c.Set(config.AccessTokenKey, accessToken)
	}
	if refreshToken != "" {
		c.Set(config.RefreshTokenKey, refreshToken)
	}
}
