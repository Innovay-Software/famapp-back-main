package handlers

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
	"github.com/innovay-software/famapp-main/config"
)

func SetAccessAndRefreshTokens(accessToken, refreshToken string, c *gin.Context) {
	if accessToken != "" {
		c.Set(config.AccessTokenKey, accessToken)
	}
	if refreshToken != "" {
		c.Set(config.RefreshTokenKey, refreshToken)
	}
}

func validateJson(req dto.ApiRequest, c *gin.Context) error {
	modelType := reflect.TypeOf(req).String()
	if reflect.ValueOf(req).Kind() != reflect.Ptr {
		return fmt.Errorf("invalid ApiRequest passed: %s", modelType)
	}

	if err := c.BindJSON(req); err != nil {
		return err
	}
	if err := validate.Struct(req); err != nil {
		return err
	}

	c.Set("requestBody", "")
	if jsonBytes, err := json.Marshal(req); err == nil {
		var jsonMap map[string]any
		// Convert to map to delete password and lockerPasscode fields
		if err := json.Unmarshal(jsonBytes, &jsonMap); err == nil {
			delete(jsonMap, "password")
			delete(jsonMap, "lockerPasscode")
			jsonBytes, err = json.Marshal(jsonMap)
			if err == nil {
				c.Set("requestBody", string(jsonBytes))
			}
		}
	}
	return nil
}

func validateUri(reqUri dto.ApiRequestUri, c *gin.Context) error {
	modelType := reflect.TypeOf(reqUri).String()
	if reflect.ValueOf(reqUri).Kind() != reflect.Ptr {
		return fmt.Errorf("invalid ApiRequest passed: %s", modelType)
	}

	if err := c.BindUri(reqUri); err != nil {
		return err
	}
	err := validate.Struct(reqUri)
	if err != nil {
		return err
	}
	return nil
}
