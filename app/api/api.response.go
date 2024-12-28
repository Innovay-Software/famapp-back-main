package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/innovay-software/famapp-main/app/dto"
	apiError "github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/innovay-software/famapp-main/config"
)

// General API response handler
func handleApiResponse(c *gin.Context, res dto.ApiResponse, err error) {
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}
	apiRespSucess(c, res)
}

func handleFileResponse(c *gin.Context, filepath string, err error) {
	if err != nil {
		apiRespFailError(c, err, nil)
		return
	}
	// returns files
	c.File(filepath)
}

// API success response handler
func apiRespSucess(c *gin.Context, data dto.ApiResponse) {
	c.IndentedJSON(http.StatusOK, apiResp(true, 0, "", false, data, c))
}

// API error response handler
func apiRespFailError(c *gin.Context, err error, data dto.ApiResponse) {
	errMessage := err.Error()
	errCode := 9999
	requiresLogin := false

	if errs, ok := err.(validator.ValidationErrors); ok {
		errMessage = utils.TranslateValidatorError(
			errs, strings.ToLower(c.GetHeader("Accept-Language")), enTrans, zhTrans,
		)
	} else if apiError, ok := err.(apiError.ApiError); ok {
		// If it's an APIError, use it's content as is
		errCode = apiError.Code
		errMessage = apiError.Message
		requiresLogin = apiError.RequiresLogin
	} else {
		// For all other errors, use the default 9999 errCode
	}

	c.IndentedJSON(http.StatusOK, apiResp(false, errCode, errMessage, requiresLogin, data, c))
}

// General API response generator
func apiResp(
	success bool, errCode int, errMsg string, requiresLogin bool, data dto.ApiResponse, c *gin.Context,
) gin.H {
	if errMsg != "" {
		utils.LogError("ApiResp error:", errCode, ", ", errMsg)
	}
	requester := ""
	userId := uint64(0)
	if user, err := getAuthenticatedUser(c); err == nil {
		userId = user.ID
		requester = user.UUID.String()
	}

	accessToken, refreshToken := "", ""
	if success {
		errCode = 0
		errMsg = ""
		if val, exists := c.Get(config.AccessTokenKey); exists {
			accessToken = val.(string)
			// data.SetAccessToken(val.(string))
		}
		if val, exists := c.Get(config.RefreshTokenKey); exists {
			refreshToken = val.(string)
			// data.SetRefreshToken(val.(string))
		}
	}

	response := map[string]any{
		"success":             success,
		"data":                data,
		"errorCode":           errCode,
		"errorMessage":        errMsg,
		"requester":           requester,
		"responseDateTime":    time.Now().UTC(),
		"hasCookie":           false,
		"ip":                  c.ClientIP(),
		"method":              c.Request.Method,
		"accessToken":         accessToken,
		"refreshToken":        refreshToken,
		"requiresLogin":       requiresLogin,
		"invalidateAllTokens": requiresLogin,
	}

	requestBodyString := ""
	requestBody, exists := c.Get("requestBody")
	if !exists {
		requestBody = ""
	} else {
		requestBodyString = requestBody.(string)
		requestLength := len(requestBodyString)
		if requestLength > 511 {
			requestBodyString = requestBodyString[:511]
		}
	}

	// Save to traffics table
	traffic := models.Traffic{
		IP:           c.ClientIP(),
		UserID:       userId,
		Requester:    requester,
		RequestURI:   c.Request.URL.RequestURI(),
		RequestBody:  requestBodyString,
		ResponseCode: "200",
		ErrorMessage: errMsg,
	}
	repositories.UtilsRepoIns.SaveTraffic(&traffic)

	return response
}
