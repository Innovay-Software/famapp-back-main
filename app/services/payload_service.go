package services

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/errors"
	"github.com/innovay-software/famapp-main/app/utils"
)

const (
	mapParamKey   = "MapParam"
	mapPayloadKey = "MapPayload"
	mapHeaderKey  = "MapHeader"
)

// Extract header data, always return the first string
func GetHeaderDataString(
	key string, defaultVal string, c *gin.Context,
) string {
	if val, exists := c.Request.Header[key]; exists {
		return val[0]
	}
	return defaultVal
}

// Extract header data, always return all strings
func GetHeaderDataSlice(
	key string, c *gin.Context,
) *[]string {
	if val, exists := c.Request.Header[key]; exists {
		return &val
	}
	return nil
}

// Set payload data to context (include request payload and url placeholders)
func SetPayloadMap(c *gin.Context) error {
	// Input data are in json format, ignore input if it is not in json
	byteList, err := c.GetRawData()
	if err != nil {
		return err
	}
	if len(byteList) == 0 {
		return nil
	}

	mapData := map[string]any{}
	err = json.Unmarshal(byteList, &mapData)
	if err != nil {
		return err
	}

	c.Set(mapPayloadKey, mapData)
	return nil
}

// Extract payload data
func getPayloadData(key string, c *gin.Context) any {
	var json map[string]any
	if raw, exists := c.Get(mapPayloadKey); exists {
		json = raw.(map[string]any)
	}

	if val, exists := json[key]; exists {
		return val
	}
	return nil
}

func GetPayloadDataMap(key string, c *gin.Context) (
	*map[string]any, error,
) {
	if item, ok := getPayloadData(key, c).(map[string]any); ok {
		return &item, nil
	}
	return nil, errors.ApiErrorParamMissing
}

func GetPayloadData[T bool | float64 | string](
	key string, defaultVal T, c *gin.Context,
) T {
	// Get payload data using key
	item := getPayloadData(key, c)
	if val, ok := item.(T); ok {
		return val
	}

	utils.Log("Get request json data failed, expected", reflect.TypeOf(defaultVal), "got", reflect.TypeOf(item))
	return defaultVal
}

func GetPayloadDataInt64(key string, defaultVal int64, c *gin.Context) int64 {
	val := GetPayloadData[float64](key, float64(defaultVal), c)
	return int64(val)
}

func GetPayloadDataSliceInt64(key string, c *gin.Context) []int64 {
	requestData := getPayloadData(key, c).([]any)
	if len(requestData) == 0 {
		return []int64{}
	}

	int64List := []int64{}
	for _, item := range requestData {
		if val, ok := item.(float64); ok {
			int64List = append(int64List, int64(val))
		}
	}
	return int64List
}

// Set path params
func SetParamMap(c *gin.Context) error {
	mapData := map[string]string{}
	for _, param := range c.Params {
		mapData[param.Key] = param.Value
	}
	c.Set(mapParamKey, mapData)
	return nil
}

// Get path param
func GetParamData(key string, defaultVal string, c *gin.Context) string {
	if raw, exists := c.Get(mapParamKey); exists {
		json := raw.(map[string]string)
		if val, exists := json[key]; exists {
			return val
		}
	}
	// if key doesn't exists, will return the zero value of string
	return defaultVal
}

func GetParamDataInt(key string, defaultVal int, c *gin.Context) int {
	val := GetParamData(key, strconv.Itoa(defaultVal), c)
	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	}
	return defaultVal
}
