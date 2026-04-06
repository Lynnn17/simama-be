package master

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"lms-be/infras"
	"lms-be/shared/failure"
	"lms-be/shared/logger"
	"lms-be/shared/model"
	"lms-be/shared/pagination"
	"lms-be/shared/repository"

	"github.com/jmoiron/sqlx"
)

var (
	DepartmentQuery = struct {
		Select                 string
		SelectDTO              string
		Insert                 string
		Update                 string
		Delete                 string
		Count                  string
		Exist                  string
		InsertBatch            string
		InsertBatchPlaceholder string
		UpsertBatchOnConflict  string
	}{
		Select: `SELECT id, name, map_coordinate_json, created_at, created_by, updated_at, updated_by, deleted_at, is_deleted FROM m_department `,
		SelectDTO: `SELECT id, name, map_coordinate_json
			FROM m_department `,
		Insert: `INSERT INTO m_department
			(name, map_coordinate_json, created_at, created_by)
			VALUES
			(:name, :map_coordinate_json, :created_at, :created_by)
			RETURNING id`,
		Update:                 `UPDATE m_department SET name = :name, map_coordinate_json = :map_coordinate_json, updated_at = :updated_at, updated_by = :updated_by, deleted_at = :deleted_at, is_deleted = :is_deleted WHERE id = :id`,
		Delete:                 `DELETE FROM m_department WHERE id = $1`,
		Count:                  `SELECT count(id) FROM m_department `,
		Exist:                  `SELECT id FROM m_department`,
		InsertBatch:            `INSERT INTO m_department (name, map_coordinate_json, created_at, created_by) VALUES `,
		InsertBatchPlaceholder: `(:name, :map_coordinate_json, :created_at, :created_by)`,
		UpsertBatchOnConflict:  `DO UPDATE SET name = EXCLUDED.name, map_coordinate_json = EXCLUDED.map_coordinate_json, updated_at = EXCLUDED.created_at, updated_by = EXCLUDED.created_by`,
	}
)

type DepartmentRepository interface {
	Create(ctx context.Context, data *Department) error
	CreateBatchTx(ctx context.Context, tx *sqlx.Tx, data []Department) error
	UpsertBatchTx(ctx context.Context, tx *sqlx.Tx, data []Department) error
	Update(ctx context.Context, data Department) error
	GetAll(ctx context.Context, req model.StandardRequest) (data []DepartmentDTO, err error)
	ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(ctx context.Context, id int) (Department Department, err error)
	DeleteByID(ctx context.Context, id int) (err error)
	ExistByField(ctx context.Context, params model.ExistByFieldParams) (bool, error)
}

type DepartmentRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideDepartmentRepositoryPostgreSQL(db *infras.PostgresqlConn) *DepartmentRepositoryPostgreSQL {
	s := new(DepartmentRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *DepartmentRepositoryPostgreSQL) Create(ctx context.Context, data *Department) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, DepartmentQuery.Insert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, data).Scan(&data.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

func (r *DepartmentRepositoryPostgreSQL) CreateBatchTx(ctx context.Context, tx *sqlx.Tx, data []Department) error {
	return r.ProcessBatchInsertTx(ctx, tx, data, false)
}

func (r *DepartmentRepositoryPostgreSQL) UpsertBatchTx(ctx context.Context, tx *sqlx.Tx, data []Department) error {
	return r.ProcessBatchInsertTx(ctx, tx, data, true)
}

func (r *DepartmentRepositoryPostgreSQL) ProcessBatchInsertTx(ctx context.Context, tx *sqlx.Tx, data []Department, upsert bool) error {
	if len(data) == 0 {
		return nil
	}

	suffix := " ON CONFLICT (name) DO NOTHING "
	if upsert {
		suffix = " ON CONFLICT (name) " + DepartmentQuery.UpsertBatchOnConflict
	}

	config := repository.BatchInsertConfig{
		BaseQuery:       DepartmentQuery.InsertBatch,
		Placeholder:     DepartmentQuery.InsertBatchPlaceholder,
		ReturningSuffix: suffix,
		MaxBatchSize:    1000,
	}

	mapFunc := func(item Department) map[string]interface{} {
		return map[string]interface{}{
			"name":                item.Name,
			"map_coordinate_json": item.MapCoordinateJSON,
			"created_at":          item.CreatedAt,
			"created_by":          item.CreatedBy,
		}
	}

	composed, err := repository.ComposeBatchInsertQuery(data, config, mapFunc)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	for _, chunk := range composed.Chunks {
		query := tx.Rebind(chunk.Query)
		_, err := tx.ExecContext(ctx, query, chunk.Args...)
		if err != nil {
			logger.ErrorWithStack(err)
			return err
		}
	}

	return nil
}

func (d *DepartmentRepositoryPostgreSQL) Update(ctx context.Context, data Department) error {
	stmt, err := d.DB.Write.PrepareNamedContext(ctx, DepartmentQuery.Update)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, data)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return err
}

func (r *DepartmentRepositoryPostgreSQL) ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false) = false")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" name ILIKE ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind("SELECT COUNT(x.id) FROM (" + DepartmentQuery.SelectDTO + searchRoleBuff.String() + ")x")

	var totalData int
	err = r.DB.Read.QueryRowContext(ctx, query, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		data.Items = make([]interface{}, 0)
		return
	}

	searchRoleBuff.WriteString(" ORDER BY " + ColumnMapDepartment[req.SortBy].(string) + " " + req.SortType)

	if !req.IgnorePaging {
		offset := (req.PageNumber - 1) * req.PageSize
		searchRoleBuff.WriteString(" LIMIT ? OFFSET ? ")
		searchParams = append(searchParams, req.PageSize)
		searchParams = append(searchParams, offset)
	}

	searchDepartmentQuery := searchRoleBuff.String()
	searchDepartmentQuery = r.DB.Read.Rebind(DepartmentQuery.SelectDTO + searchDepartmentQuery)

	rows, err := r.DB.Read.QueryxContext(ctx, searchDepartmentQuery, searchParams...)

	if err != nil {
		return
	}

	for rows.Next() {
		var departmentDTO DepartmentDTO
		err = rows.StructScan(&departmentDTO)
		if err != nil {
			return
		}
		data.Items = append(data.Items, departmentDTO)
	}

	if req.IgnorePaging {
		data.Meta = pagination.CreateMeta(totalData, totalData, 1)
	} else {
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	}

	return
}

func (r *DepartmentRepositoryPostgreSQL) GetAll(ctx context.Context, req model.StandardRequest) (data []DepartmentDTO, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false) = false")
	searchRoleBuff.WriteString(" ORDER BY created_at DESC ")
	query := r.DB.Read.Rebind(DepartmentQuery.SelectDTO + searchRoleBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, searchParams...)

	if err == sql.ErrNoRows {
		_ = failure.NotFound("Department not found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var departmentDTO DepartmentDTO
		err = rows.StructScan(&departmentDTO)

		if err != nil {
			return
		}
		data = append(data, departmentDTO)
	}

	return
}

func (r *DepartmentRepositoryPostgreSQL) ResolveByID(ctx context.Context, id int) (department Department, err error) {
	err = r.DB.Read.GetContext(ctx, &department, DepartmentQuery.Select+" WHERE id=$1", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *DepartmentRepositoryPostgreSQL) DeleteByID(ctx context.Context, id int) (err error) {
	_, err = r.DB.Write.ExecContext(ctx, DepartmentQuery.Delete, id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *DepartmentRepositoryPostgreSQL) ExistByField(ctx context.Context, params model.ExistByFieldParams) (bool, error) {
	var idResult int
	where := fmt.Sprintf(" WHERE UPPER(%s) = UPPER($1) AND coalesce(is_deleted, false) = false", params.Field)
	if params.ExcludeID != 0 {
		where += fmt.Sprintf(" AND id <> %d", params.ExcludeID)
	}

	err := r.DB.Read.GetContext(ctx, &idResult, DepartmentQuery.Exist+where, params.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.ErrorWithStack(err)
		return false, err
	}

	return true, nil
}
