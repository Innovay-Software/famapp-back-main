package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validatorV10 "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	apiErrors "github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/middlewares"
	"github.com/innovay-software/famapp-main/config"
	ginMiddleware "github.com/oapi-codegen/gin-middleware"
)

var (
	uni      *ut.UniversalTranslator
	validate *validatorV10.Validate
	zhTrans  ut.Translator
	enTrans  ut.Translator
)

func RegisterRoutes(r *gin.Engine, isProd bool) error {
	// Init validator
	initValidator()

	// Default panic recovery
	r.Use(gin.CustomRecovery(apiPanicRecoverHandler))

	// Default 404
	r.NoRoute(api404Handler)

	// Set CORS
	if !isProd {
		r.Use(middlewares.CorsSettings())
	}

	// Integrate OpenAPI authentication func
	// Wheter or not an endpoint requires authorization is
	// controled in the OpenAPI specifications
	validatorOptions := &ginMiddleware.Options{}
	validatorOptions.Options.AuthenticationFunc = jwtAuthorizationFunc

	// Get swagger instance
	swagger, err := GetSwagger()
	if err != nil {
		return fmt.Errorf("error loading swagger specs: %s", err)
	}
	swagger.Servers = nil
	r.Use(ginMiddleware.OapiRequestValidatorWithOptions(
		swagger, validatorOptions,
	))

	// Set JWT authentication middleware
	r.Use(middlewares.JwtAuthentication())

	// Register handlers
	v2 := r.Group("/")
	RegisterHandlersWithOptions(
		v2, &ApiServerInterfaceImpl{}, GinServerOptions{},
	)

	return nil
}

func initValidator() {
	validate = validatorV10.New()
	en := en.New()
	zh := zh.New()
	uni = ut.New(en, zh)
	zhTrans, _ = uni.GetTranslator("zh")
	zh_translations.RegisterDefaultTranslations(validate, zhTrans)
	enTrans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, enTrans)
}

func jwtAuthorizationFunc(c context.Context, input *openapi3filter.AuthenticationInput) error {
	// Set a flag in header so the authentication middleware will pick it up and do user authentication
	input.RequestValidationInput.Request.Header.Set(config.RequireJwtHeaderKey, "user")
	return nil
}

// Handles middle panics, prints out the error in JSON format and aborts.
func apiPanicRecoverHandler(c *gin.Context, err any) {
	if apiError, ok := err.(apiErrors.ApiError); ok {
		apiRespFailError(c, apiError, nil)
	} else {
		errMessage := ""
		if str, ok := err.(string); ok {
			errMessage = str
		}
		if errErr, ok := err.(error); ok {
			errMessage = errErr.Error()
		}
		apiRespFailError(c, apiErrors.ApiError{Code: -1, Message: errMessage}, nil)
	}

	c.Abort()
}

// 404 not found handler
func api404Handler(c *gin.Context) {
	c.IndentedJSON(
		http.StatusNotFound,
		apiResp(false, apiErrors.ApiError404.Code, apiErrors.ApiError404.Message,
			apiErrors.ApiError404.RequiresLogin, nil, c,
		),
	)
}
