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
	CompanyQuery = struct {
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
		Select: `SELECT id, name, is_registered_partner, pic_contact, created_at, created_by, updated_at, updated_by, deleted_at, is_deleted FROM m_company `,
		SelectDTO: `SELECT id, name, is_registered_partner, pic_contact
			FROM m_company `,
		Insert: `INSERT INTO m_company
			(name, is_registered_partner, pic_contact, created_at, created_by)
			VALUES
			(:name, :is_registered_partner, :pic_contact, :created_at, :created_by)
			RETURNING id`,
		Update:                 `UPDATE m_company SET name = :name, is_registered_partner = :is_registered_partner, pic_contact = :pic_contact,  updated_at = :updated_at, updated_by = :updated_by, deleted_at = :deleted_at, is_deleted = :is_deleted WHERE id = :id`,
		Delete:                 `DELETE FROM m_company WHERE id = $1`,
		Count:                  `SELECT count(id) FROM m_company `,
		Exist:                  `SELECT id FROM m_company`,
		InsertBatch:            `INSERT INTO m_company (name, is_registered_partner, pic_contact, created_at, created_by) VALUES `,
		InsertBatchPlaceholder: `(:name, :is_registered_partner, :pic_contact, :created_at, :created_by)`,
		UpsertBatchOnConflict:  `DO UPDATE SET name = EXCLUDED.name, is_registered_partner = EXCLUDED.is_registered_partner, pic_contact = EXCLUDED.pic_contact, updated_at = EXCLUDED.created_at, updated_by = EXCLUDED.created_by`,
	}
)

type CompanyRepository interface {
	Create(ctx context.Context, data *Company) error
	CreateBatchTx(ctx context.Context, tx *sqlx.Tx, data []Company) error
	UpsertBatchTx(ctx context.Context, tx *sqlx.Tx, data []Company) error
	Update(ctx context.Context, data Company) error
	GetAll(ctx context.Context, req model.StandardRequest) (data []CompanyDTO, err error)
	ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(ctx context.Context, id int) (Company Company, err error)
	DeleteByID(ctx context.Context, id int) (err error)
	ExistByField(ctx context.Context, params model.ExistByFieldParams) (bool, error)
}

type CompanyRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideCompanyRepositoryPostgreSQL(db *infras.PostgresqlConn) *CompanyRepositoryPostgreSQL {
	s := new(CompanyRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *CompanyRepositoryPostgreSQL) Create(ctx context.Context, data *Company) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, CompanyQuery.Insert)
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

func (r *CompanyRepositoryPostgreSQL) CreateBatchTx(ctx context.Context, tx *sqlx.Tx, data []Company) error {
	return r.ProcessBatchInsertTx(ctx, tx, data, false)
}

func (r *CompanyRepositoryPostgreSQL) UpsertBatchTx(ctx context.Context, tx *sqlx.Tx, data []Company) error {
	return r.ProcessBatchInsertTx(ctx, tx, data, true)
}

func (r *CompanyRepositoryPostgreSQL) ProcessBatchInsertTx(ctx context.Context, tx *sqlx.Tx, data []Company, upsert bool) error {
	if len(data) == 0 {
		return nil
	}

	suffix := " ON CONFLICT (name) DO NOTHING "
	if upsert {
		suffix = " ON CONFLICT (name) " + CompanyQuery.UpsertBatchOnConflict
	}

	config := repository.BatchInsertConfig{
		BaseQuery:       CompanyQuery.InsertBatch,
		Placeholder:     CompanyQuery.InsertBatchPlaceholder,
		ReturningSuffix: suffix,
		MaxBatchSize:    1000,
	}

	mapFunc := func(item Company) map[string]interface{} {
		return map[string]interface{}{
			"name":                  item.Name,
			"is_registered_partner": item.IsRegisteredPartner,
			"pic_contact":           item.PICContact,
			"created_at":            item.CreatedAt,
			"created_by":            item.CreatedBy,
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

func (c *CompanyRepositoryPostgreSQL) Update(ctx context.Context, data Company) error {
	stmt, err := c.DB.Write.PrepareNamedContext(ctx, CompanyQuery.Update)
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

func (r *CompanyRepositoryPostgreSQL) ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false) = false")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(name, pic_contact) ILIKE ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind("SELECT COUNT(x.id) FROM (" + CompanyQuery.SelectDTO + searchRoleBuff.String() + ")x")

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

	searchRoleBuff.WriteString(" ORDER BY " + ColumnMapCompany[req.SortBy].(string) + " " + req.SortType)

	if !req.IgnorePaging {
		offset := (req.PageNumber - 1) * req.PageSize
		searchRoleBuff.WriteString(" LIMIT ? OFFSET ? ")
		searchParams = append(searchParams, req.PageSize)
		searchParams = append(searchParams, offset)
	}

	searchCompanyQuery := searchRoleBuff.String()
	searchCompanyQuery = r.DB.Read.Rebind(CompanyQuery.SelectDTO + searchCompanyQuery)

	rows, err := r.DB.Read.QueryxContext(ctx, searchCompanyQuery, searchParams...)

	if err != nil {
		return
	}

	for rows.Next() {
		var companyDTO CompanyDTO
		err = rows.StructScan(&companyDTO)
		if err != nil {
			return
		}
		data.Items = append(data.Items, companyDTO)
	}

	if req.IgnorePaging {
		data.Meta = pagination.CreateMeta(totalData, totalData, 1)
	} else {
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	}

	return
}

func (r *CompanyRepositoryPostgreSQL) GetAll(ctx context.Context, req model.StandardRequest) (data []CompanyDTO, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false) = false")
	searchRoleBuff.WriteString(" ORDER BY created_at DESC ")
	query := r.DB.Read.Rebind(CompanyQuery.SelectDTO + searchRoleBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, searchParams...)

	if err == sql.ErrNoRows {
		_ = failure.NotFound("Company not found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var companyDTO CompanyDTO
		err = rows.StructScan(&companyDTO)

		if err != nil {
			return
		}
		data = append(data, companyDTO)
	}

	return
}

func (r *CompanyRepositoryPostgreSQL) ResolveByID(ctx context.Context, id int) (company Company, err error) {
	err = r.DB.Read.GetContext(ctx, &company, CompanyQuery.Select+" WHERE id=$1", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *CompanyRepositoryPostgreSQL) DeleteByID(ctx context.Context, id int) (err error) {
	_, err = r.DB.Write.ExecContext(ctx, CompanyQuery.Delete, id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *CompanyRepositoryPostgreSQL) ExistByField(ctx context.Context, params model.ExistByFieldParams) (bool, error) {
	var idResult int
	where := fmt.Sprintf(" WHERE UPPER(%s) = UPPER($1) AND coalesce(is_deleted, false) = false", params.Field)
	if params.ExcludeID != 0 {
		where += fmt.Sprintf(" AND id <> %d", params.ExcludeID)
	}

	err := r.DB.Read.GetContext(ctx, &idResult, CompanyQuery.Exist+where, params.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.ErrorWithStack(err)
		return false, err
	}

	return true, nil
}
