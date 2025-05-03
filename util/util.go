package util

import (
	"encoding/json"
	"net/http"

	"github.com/glossd/fetch"
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

type JsonLower struct {
	Data any // 要輸出的資料
}

func (r JsonLower) Render(w http.ResponseWriter) error {
	jsonData, err := fetch.Marshal(r.Data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(jsonData))
	return err
}

func (r JsonLower) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func JsonL(data any) JsonLower {
	return JsonLower{Data: data}
}
