package integrationTests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
	"github.com/innovay-software/famapp-main/app/utils"
)

func logWorkflowSuccess(content string) {
	utils.LogSuccess("\n***" + content + "***\n")
}

// Initiate a Post Request
func postRequest(
	router *gin.Engine, uri string, headers *map[string]string,
	requestBody any, responseModel dto.ApiResponse,
) error {
	return initHttpRequest(router, "POST", uri, headers, requestBody, responseModel)
}

// Initiate a Get request
func getRequest(
	router *gin.Engine, uri string, headers *map[string]string,
	responseModel dto.ApiResponse,
) error {
	return initHttpRequest(router, "GET", uri, headers, "", responseModel)
}

// Initiate a request based on provided method
func initHttpRequest(
	router *gin.Engine, method, uri string, headers *map[string]string,
	requestBody any, responseModel dto.ApiResponse,
) error {
	bodyJson, _ := json.Marshal(requestBody)
	bodyReader := bytes.NewReader(bodyJson)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, uri, bodyReader)
	defaultHeaders := generateDefaultHeadersForAPIRequests("")
	if headers != nil {
		for k, v := range *headers {
			(*defaultHeaders)[k] = v
		}
	}
	for k, v := range *defaultHeaders {
		req.Header.Add(k, v)
	}

	router.ServeHTTP(w, req)
	return processApiResponse(w.Body.Bytes(), responseModel)
}

// Post request for binary data
func postRequestBinaryPayload(
	router *gin.Engine, uri string, headers *map[string]string,
	binaryPayload *[]byte, responseModel dto.ApiResponse,
) error {
	bodyReader := bytes.NewReader(*binaryPayload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", uri, bodyReader)
	defaultHeaders := generateDefaultHeadersForAPIRequests("")
	if headers != nil {
		for k, v := range *headers {
			(*defaultHeaders)[k] = v
		}
	}
	for k, v := range *defaultHeaders {
		req.Header.Add(k, v)
	}

	router.ServeHTTP(w, req)
	return processApiResponse(w.Body.Bytes(), responseModel)
}

// Processes the api response.
// Validates if it's 200, then if it was successful, then unmarshals to target ApiResponse model
func processApiResponse(
	bytes []byte, responseModel dto.ApiResponse,
) error {
	jsonMap := map[string]any{}
	if err := json.Unmarshal(bytes, &jsonMap); err != nil {
		return err
	}
	if jsonMap["success"] == nil {
		utils.LogError("Invalid api response")
		utils.LogError(string(bytes))
		return errors.New("Invalid api response format")
	}

	if jsonMap["accessToken"] != "" {
		responseModel.SetAccessToken(jsonMap["accessToken"].(string))
	}
	if jsonMap["refreshToken"] != "" {
		responseModel.SetRefreshToken(jsonMap["refreshToken"].(string))
	}

	if jsonMap["success"] != true || jsonMap["errorMessage"] != "" {
		utils.LogError("Test: API call failed:", jsonMap["errorMessage"])
		errMessage := "API call failed: " + string(bytes)
		if jsonMap["errorMessage"] != nil {
			errMessage = errMessage + jsonMap["errorMessage"].(string)
		}
		return errors.New(errMessage)
	}

	// Extract the data field of the response, and unmarshal into the responseModel
	dataJsonMap, ok := jsonMap["data"].(map[string]any)
	if !ok {
		return errors.New("Unabled to extract data from response")
	}
	dataJsonString, err := json.Marshal(dataJsonMap)
	if err != nil {
		return errors.New("Response marshal error")
	}
	err2 := json.Unmarshal(dataJsonString, responseModel)
	if err2 != nil {
		return errors.New("Response marshal.unmarshal error")
	}

	return nil
}

// Default headers for API requests
func generateDefaultHeadersForAPIRequests(token string) *map[string]string {
	return &map[string]string{
		// "Integration-Testing": "1",
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}
}
