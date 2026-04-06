package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

const RECORD_NOT_FOUND = "record not found"

// StandardRequest is a standard query string request
type StandardRequest struct {
	Keyword      string `json:"q" validate:"omitempty"`
	StartDate    string `json:"startDate" validate:"omitempty"`
	EndDate      string `json:"endDate" validate:"omitempty"`
	PageNumber   int    `json:"pageNumber" validate:"omitempty,gte=0"`
	PageSize     int    `json:"pageSize" validate:"omitempty,gte=0"`
	SortBy       string `json:"sortBy" validate:"required"`
	SortType     string `json:"sortType" validate:"required,oneof=asc ASC desc DESC"`
	Status       string `json:"status" validate:"omitempty"`
	IgnorePaging bool   `json:"ignorePaging" validate:"omitempty"`
	RoleID       string `json:"roleId" validate:"omitempty"`
	PersonID     string `json:"personId" validate:"omitempty"`
}
type StandardModel struct {
	ID     string `json:"id" db:"id"`
	Nama   string `json:"nama" db:"nama"`
	Alamat string `json:"alamat" db:"alamat"`
}
type StandardRequestUser struct {
	Keyword    string `json:"q" validate:"omitempty"`
	StartDate  string `json:"startDate" validate:"omitempty"`
	EndDate    string `json:"endDate" validate:"omitempty"`
	PageNumber int    `json:"pageNumber" validate:"omitempty,gte=0"`
	PageSize   int    `json:"pageSize" validate:"omitempty,gte=0"`
	SortBy     string `json:"sortBy" validate:"required"`
	SortType   string `json:"sortType" validate:"required,oneof=asc ASC desc DESC"`
	Active     bool   `json:"active" validate:"omitempty"`
	RoleID     string `json:"roleId" validate:"omitempty"`
	PersonID   string `json:"personId" validate:"omitempty"`
}

type StandardRequestMenu struct {
	Keyword    string `json:"q" validate:"omitempty"`
	PageNumber int    `json:"pageNumber" validate:"omitempty,gte=0"`
	PageSize   int    `json:"pageSize" validate:"omitempty,gte=0"`
	SortBy     string `json:"sortBy" validate:"required"`
	SortType   string `json:"sortType" validate:"required,oneof=asc ASC desc DESC"`
	App        string `json:"app" validate:"omitempty"`
}

type ExistByFieldParams struct {
	Field         string
	Value         string
	ExcludeID     int    // ID to exclude from check (for updates)
	ExcludeIDStr  string // ID to exclude from check (for updates)
	InstitutionID int
}

// JSONRaw ...
type JSONRaw json.RawMessage

// Value ...
func (j JSONRaw) Value() (driver.Value, error) {
	byteArr := []byte(j)

	return driver.Value(byteArr), nil
}

// Scan ...
func (j *JSONRaw) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	err := json.Unmarshal(asBytes, &j)
	if err != nil {
		return error(errors.New("Scan could not unmarshal to []string"))
	}

	return nil
}

// MarshalJSON ...
func (j *JSONRaw) MarshalJSON() ([]byte, error) {
	return *j, nil
}

// UnmarshalJSON ...
func (j *JSONRaw) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}
