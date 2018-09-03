package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

type MetadataMap map[string]interface{}

func (m MetadataMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m MetadataMap) Scan(src interface{}) error {
	v := reflect.ValueOf(src)
	if !v.IsValid() || v.IsNil() {
		return nil
	}
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, &m)
	}
	return fmt.Errorf("Could not not decode type %T -> %T", src, m)
}
