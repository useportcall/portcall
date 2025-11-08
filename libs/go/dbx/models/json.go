package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type InvoiceConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m InvoiceConfig) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *InvoiceConfig) Scan(src any) error {
	var data []byte

	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}

	return json.Unmarshal(data, m)
}

type Tier struct {
	Start  int `json:"start"`
	End    int `json:"end"`
	Amount int `json:"unit_amount"`
}
