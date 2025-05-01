package util

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func GetDefault[T any](a *T) T {
	if a == nil {
		var t T
		return t
	}
	return *a
}

func ToDataTypeJSON(arr []string) (datatypes.JSON, error) {
	json, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(json), nil
}
