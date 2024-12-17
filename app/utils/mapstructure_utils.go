package utils

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func PopulateModelFromMap[T any](modelInstance *T, data *map[string]any) error {
	err := mapstructure.Decode(data, modelInstance)
	if err != nil {
		Log(reflect.TypeOf((*data)["metadata"]))
		LogError("Populate model from map failed:", err)
		return err
	}
	return nil
}
