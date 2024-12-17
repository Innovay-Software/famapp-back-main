package tests

import (
	"reflect"
	"slices"
	"testing"

	"github.com/innovay-software/famapp-main/app/utils"
)

func AssertNil(t *testing.T, param any) {
	if param == nil {
		return
	}
	if !reflect.ValueOf(param).IsNil() {
		utils.LogError("AssertNil failed")
		utils.LogError(param)
		t.Errorf("Error")
	}
}

func AssertNotNil(t *testing.T, param any) {
	if param == nil || reflect.ValueOf(param).IsNil() {
		utils.LogError("AssertNotNil failed")
		t.Errorf("Error")
	}
}

func AssertEqual(t *testing.T, actualResult any, expectedResult any) {
	if actualResult != expectedResult {
		utils.LogError("AssertEqual Error, expected", expectedResult, "got", actualResult)
		t.Errorf("Error")
	}
}

func AssertNotEqual(t *testing.T, actualResult any, expectedResult any) {
	if actualResult == expectedResult {
		utils.LogError("AssertNotEqual Error, expected to not equal to:", expectedResult, "-")
		t.Errorf("Error")
	}
}

func AssertIn(t *testing.T, actualResult any, candidates []any) {
	if !slices.Contains(candidates, actualResult) {
		utils.LogError("AssertIn Error, ", actualResult, "is not in", candidates)
		t.Errorf("Error")
	}
}
