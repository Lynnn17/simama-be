package model

// ImportResult represents a generic result of an import operation.
// This can be used by any domain that performs batch imports.
type ImportResult struct {
	TotalRows    int              `json:"totalRows"`
	SuccessCount int              `json:"successCount"`
	UpdateCount  int              `json:"updateCount"`
	SkipCount    int              `json:"skipCount"`
	ErrorCount   int              `json:"errorCount"`
	Errors       []ImportRowError `json:"errors,omitempty"`
}

// ImportRowError represents an error for a specific row during import.
type ImportRowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

// PreviewResult represents a generic preview result before import.
// The Data field uses interface{} to support different row types.
type PreviewResult struct {
	TotalRows  int              `json:"totalRows"`
	ValidCount int              `json:"validCount"`
	ErrorCount int              `json:"errorCount"`
	Data       interface{}      `json:"data"`
	Errors     []ImportRowError `json:"errors,omitempty"`
}
