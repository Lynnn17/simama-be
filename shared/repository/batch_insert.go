package repository

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

// BatchInsertConfig contains configuration for batch insert
type BatchInsertConfig struct {
	BaseQuery       string // e.g. `INSERT INTO table_name (col1, col2) VALUES `
	Placeholder     string // e.g. `(:col1, :col2)`
	ReturningSuffix string // e.g. ` RETURNING id` (optional)
	MaxBatchSize    int    // e.g. 1000 (default if 0)
}

// BatchInsertResult contains the result of batch insert query composition
type BatchInsertResult struct {
	Query  string
	Args   []interface{}
	Chunks []BatchInsertChunk
}

// BatchInsertChunk represents a single chunk of batch insert
type BatchInsertChunk struct {
	StartIndex int
	EndIndex   int
	Query      string
	Args       []interface{}
}

// MapFunc is a function type that converts a struct to a map for named parameters
type MapFunc[T any] func(item T) map[string]interface{}

// ComposeBatchInsertQuery builds batch insert queries with chunking support
// T is the type of the data slice
// mapFunc converts each item to a map[string]interface{} for sqlx.Named
func ComposeBatchInsertQuery[T any](data []T, config BatchInsertConfig, mapFunc MapFunc[T]) (result BatchInsertResult, err error) {
	if len(data) == 0 {
		return result, nil
	}

	// Default max batch size
	maxBatchSize := config.MaxBatchSize
	if maxBatchSize <= 0 {
		maxBatchSize = 1000
	}

	// Process data in chunks
	for chunkStart := 0; chunkStart < len(data); chunkStart += maxBatchSize {
		chunkEnd := chunkStart + maxBatchSize
		if chunkEnd > len(data) {
			chunkEnd = len(data)
		}

		chunk := data[chunkStart:chunkEnd]
		values := []string{}
		var args []interface{}

		for _, item := range chunk {
			param := mapFunc(item)
			q, paramArgs, err := sqlx.Named(config.Placeholder, param)
			if err != nil {
				return result, err
			}
			values = append(values, q)
			args = append(args, paramArgs...)
		}

		query := config.BaseQuery + strings.Join(values, ", ") + config.ReturningSuffix

		result.Chunks = append(result.Chunks, BatchInsertChunk{
			StartIndex: chunkStart,
			EndIndex:   chunkEnd,
			Query:      query,
			Args:       args,
		})
	}

	return result, nil
}

// ComposeSingleBatchInsertQuery builds a single batch insert query (no chunking)
// Use this when you know the data size is small
func ComposeSingleBatchInsertQuery[T any](data []T, config BatchInsertConfig, mapFunc MapFunc[T]) (query string, args []interface{}, err error) {
	if len(data) == 0 {
		return "", nil, nil
	}

	values := []string{}
	for _, item := range data {
		param := mapFunc(item)
		q, paramArgs, err := sqlx.Named(config.Placeholder, param)
		if err != nil {
			return "", nil, err
		}
		values = append(values, q)
		args = append(args, paramArgs...)
	}

	query = config.BaseQuery + strings.Join(values, ", ") + config.ReturningSuffix
	return query, args, nil
}
