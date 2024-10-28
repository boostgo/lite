package sql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONRaw json.RawMessage

func (jr *JSONRaw) Value() (driver.Value, error) {
	byteArr := []byte(*jr)

	return driver.Value(byteArr), nil
}

func (jr *JSONRaw) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}

	if err := json.Unmarshal(asBytes, &jr); err != nil {
		return errors.New("Scan could not unmarshal to []string")
	}

	return nil
}

func (jr *JSONRaw) MarshalJSON() ([]byte, error) {
	return *jr, nil
}

func (jr *JSONRaw) UnmarshalJSON(data []byte) error {
	if jr == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	*jr = append((*jr)[0:0], data...)
	return nil
}
