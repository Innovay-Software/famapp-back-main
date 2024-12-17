package dto

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/innovay-software/famapp-main/app/utils"
	"gorm.io/gorm/schema"
)

type UnimplementedResponse struct {
	ApiResponseBase `json:",squash"`
	Comment         string `json:"comment"`
}

// Sync data from ApiRequest to schema.Tabler
// First converts to map, then converts to schema.Tabler
func SyncApiRequestToModel(
	req ApiRequest, model schema.Tabler, omitFieldNames []string,
) error {
	if reflect.ValueOf(req).Kind() != reflect.Ptr ||
		reflect.ValueOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("params should be pointers")
	}

	var mapData map[string]any
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonBytes, &mapData); err != nil {
		return err
	}
	log.Println("SyncApiRequestToModel:")
	log.Println(req)
	log.Println(mapData)

	for _, fieldName := range omitFieldNames {
		delete(mapData, fieldName)
	}

	return utils.PopulateModelFromMap(&model, &mapData)
}
