package handlers

// import (
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator/v10"
// 	"github.com/innovay-software/famapp-main/app/dto"
// 	"github.com/innovay-software/famapp-main/app/errors"
// 	"github.com/innovay-software/famapp-main/app/models"
// 	"github.com/innovay-software/famapp-main/app/repositories"
// 	"github.com/innovay-software/famapp-main/app/utils"
// )

// // General API success response
// func apiRespSucess(data dto.ApiResponse, c *gin.Context) bool {
// 	c.IndentedJSON(http.StatusOK, apiResp(true, "", data, c))
// 	return true
// }

// // General error response
// func apiRespFailError(err error, data dto.ApiResponse, c *gin.Context) bool {
// 	errMessage := err.Error()

// 	if errs, ok := err.(validator.ValidationErrors); ok {
// 		// Translate messages according to "Accept-Language" heander
// 		language := strings.ToLower(c.GetHeader("Accept-Language"))
// 		trans := enTrans
// 		if strings.HasPrefix(language, "zh") {
// 			trans = zhTrans
// 		}

// 		translatedMessages := []string{}
// 		for _, e := range errs {
// 			// can translate each error one at a time.
// 			translatedMessages = append(translatedMessages, e.Translate(trans))
// 		}

// 		errMessage = strings.Join(translatedMessages, ". ")
// 		utils.LogError("ApiRespFailError:", errMessage)
// 	}
// 	c.IndentedJSON(http.StatusOK, apiResp(false, errMessage, data, c))
// 	return false
// }

// // General API response generator
// func apiResp(success bool, errMsg string, data dto.ApiResponse, c *gin.Context) gin.H {
// 	requester := ""
// 	userId := int64(0)
// 	if user, err := getAuthenticatedUser(c); err == nil {
// 		userId = user.ID
// 		requester = user.UUID.String()
// 	}

// 	if success {
// 		errMsg = ""
// 		if val, exists := c.Get("accessToken"); exists {
// 			// log.Println("Has access token")
// 			data.SetAccessToken(val.(string))
// 		}
// 		if val, exists := c.Get("refreshToken"); exists {
// 			// log.Println("Has refresh token")
// 			data.SetRefreshToken(val.(string))
// 		}
// 		// log.Println("After:", data.GetAccessToken(), "and", data.GetRefreshToken(), "-")
// 	}

// 	response := map[string]any{
// 		"success":          success,
// 		"data":             data,
// 		"errorMessage":     errMsg,
// 		"requester":        requester,
// 		"responseDateTime": time.Now().UTC(),
// 		"hasCookie":        false,
// 		"ip":               c.ClientIP(),
// 		"method":           c.Request.Method,
// 	}

// 	requestBodyString := ""
// 	requestBody, exists := c.Get("requestBody")
// 	if !exists {
// 		requestBody = ""
// 	} else {
// 		requestBodyString = requestBody.(string)
// 		requestLength := len(requestBodyString)
// 		if requestLength > 511 {
// 			requestBodyString = requestBodyString[:511]
// 		}
// 	}

// 	// Save to traffics table
// 	traffic := models.Traffic{
// 		IP:           c.ClientIP(),
// 		UserID:       userId,
// 		Requester:    requester,
// 		RequestURI:   c.Request.URL.RequestURI(),
// 		RequestBody:  requestBodyString,
// 		ResponseCode: "200",
// 		ErrorMessage: errMsg,
// 	}
// 	repositories.SaveDbModel(&traffic)

// 	return response
// }

// // Handles middle panics, prints out the error in JSON format and aborts.
// func ApiPanicRecoverHandler(c *gin.Context, err any) {
// 	if apiError, ok := err.(errors.ApiError); ok {
// 		apiRespFailError(apiError, nil, c)
// 	} else {
// 		errMessage := ""
// 		if str, ok := err.(string); ok {
// 			errMessage = str
// 		}
// 		if errErr, ok := err.(error); ok {
// 			errMessage = errErr.Error()
// 		}
// 		apiRespFailError(errors.ApiError{Code: -1, Message: errMessage}, nil, c)
// 	}

// 	c.Abort()
// }

// // 404 not found handler
// func Api404Handler(c *gin.Context) {
// 	c.IndentedJSON(
// 		http.StatusNotFound,
// 		apiResp(false, errors.ApiError404.Message, nil, c),
// 	)
// }
