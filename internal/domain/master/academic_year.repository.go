package master

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"lms-be/infras"
	"lms-be/shared/failure"
	"lms-be/shared/logger"
	"lms-be/shared/model"
	"lms-be/shared/pagination"
	"lms-be/shared/repository"
)

var (
	AcademicYearQuery = struct {
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
		Select: `SELECT id, code, name, institution_id, semester_type, is_active, created_at, created_by, updated_at, updated_by, deleted_at, is_deleted FROM m_academic_year `,
		SelectDTO: `SELECT ay.id, ay.code, ay.name, ay.institution_id, ay.semester_type, ay.is_active, ay.created_at, ay.created_by, ay.updated_at, ay.updated_by, ay.deleted_at, ay.is_deleted, i.name as institution_name
			FROM m_academic_year ay
			LEFT JOIN m_institution i ON ay.institution_id=i.id`,
		Insert: `INSERT INTO m_academic_year
			(code, name, institution_id, semester_type, is_active, created_at, created_by, updated_at, updated_by)
			VALUES
			(:code, :name, :institution_id, :semester_type, :is_active, :created_at, :created_by, :updated_at, :updated_by) RETURNING id`,
		Update:                 `UPDATE m_academic_year SET code = :code, name = :name, institution_id = :institution_id, semester_type = :semester_type, is_active = :is_active, updated_at = :updated_at, updated_by = :updated_by, deleted_at = :deleted_at, is_deleted = :is_deleted WHERE id = :id`,
		Delete:                 `DELETE FROM m_academic_year WHERE id = $1`,
		Count:                  `SELECT count(id) FROM m_academic_year `,
		Exist:                  `SELECT id FROM m_academic_year`,
		InsertBatch:            `INSERT INTO m_academic_year (code, name, institution_id, semester_type, is_active, created_at, created_by) VALUES `,
		InsertBatchPlaceholder: `(:code, :name, :institution_id, :semester_type, :is_active, :created_at, :created_by)`,
		UpsertBatchOnConflict:  `DO UPDATE SET name = EXCLUDED.name, semester_type = EXCLUDED.semester_type, is_active = EXCLUDED.is_active, updated_at = EXCLUDED.created_at, updated_by = EXCLUDED.created_by`,
	}
)

type AcademicYearRepository interface {
	Create(ctx context.Context, data *AcademicYear) error
	CreateBatchTx(ctx context.Context, tx *sqlx.Tx, data []AcademicYear) error
	UpsertBatchTx(ctx context.Context, tx *sqlx.Tx, data []AcademicYear) error
	Update(ctx context.Context, data AcademicYear) error
	GetAll(ctx context.Context, req model.StandardRequest) (data []AcademicYearDTO, err error)
	ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(ctx context.Context, id int) (AcademicYear AcademicYear, err error)
	ResolveByCode(ctx context.Context, code string, institutionID int) (AcademicYear AcademicYear, err error)
	DeleteByID(ctx context.Context, id int) (err error)
	ExistByField(ctx context.Context, params model.ExistByFieldParams) (bool, error)
}

type AcademicYearRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideAcademicYearRepositoryPostgreSQL(db *infras.PostgresqlConn) *AcademicYearRepositoryPostgreSQL {
	s := new(AcademicYearRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *AcademicYearRepositoryPostgreSQL) Create(ctx context.Context, data *AcademicYear) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, AcademicYearQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
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

func (r *AcademicYearRepositoryPostgreSQL) CreateBatchTx(ctx context.Context, tx *sqlx.Tx, data []AcademicYear) error {
	return r.ProcessBatchInsertTx(ctx, tx, data, false)
}

func (r *AcademicYearRepositoryPostgreSQL) UpsertBatchTx(ctx context.Context, tx *sqlx.Tx, data []AcademicYear) error {
	return r.ProcessBatchInsertTx(ctx, tx, data, true)
}

func (r *AcademicYearRepositoryPostgreSQL) ProcessBatchInsertTx(ctx context.Context, tx *sqlx.Tx, data []AcademicYear, upsert bool) error {
	if len(data) == 0 {
		return nil
	}

	// Build config with optional ON CONFLICT suffix
	suffix := " ON CONFLICT (code, institution_id) DO NOTHING "
	if upsert {
		suffix = " ON CONFLICT (code, institution_id) " + AcademicYearQuery.UpsertBatchOnConflict
	}

	config := repository.BatchInsertConfig{
		BaseQuery:       AcademicYearQuery.InsertBatch,
		Placeholder:     AcademicYearQuery.InsertBatchPlaceholder,
		ReturningSuffix: suffix,
		MaxBatchSize:    1000,
	}

	mapFunc := func(item AcademicYear) map[string]interface{} {
		return map[string]interface{}{
			"code":           item.Code,
			"name":           item.Name,
			"institution_id": item.InstitutionID,
			"semester_type":  item.SemesterType,
			"is_active":      item.IsActive,
			"created_at":     item.CreatedAt,
			"created_by":     item.CreatedBy,
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

func (r *AcademicYearRepositoryPostgreSQL) Update(ctx context.Context, data AcademicYear) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, AcademicYearQuery.Update)
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

func (r *AcademicYearRepositoryPostgreSQL) ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(ay.is_deleted, false) = false")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(ay.name, ay.code, ay.semester_type, i.name) ilike ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind("SELECT COUNT(x.id) FROM (" + AcademicYearQuery.SelectDTO + searchRoleBuff.String() + ")x")

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

	searchRoleBuff.WriteString(" ORDER BY " + ColumnMapAcademicYear[req.SortBy].(string) + " " + req.SortType)
	if !req.IgnorePaging {
		offset := (req.PageNumber - 1) * req.PageSize
		searchRoleBuff.WriteString(" LIMIT ? OFFSET ? ")
		searchParams = append(searchParams, req.PageSize)
		searchParams = append(searchParams, offset)
	}

	searchAcademicYearQuery := searchRoleBuff.String()
	searchAcademicYearQuery = r.DB.Read.Rebind(AcademicYearQuery.SelectDTO + searchAcademicYearQuery)
	rows, err := r.DB.Read.QueryxContext(ctx, searchAcademicYearQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var academicYear AcademicYearDTO
		err = rows.StructScan(&academicYear)
		if err != nil {
			return
		}
		data.Items = append(data.Items, academicYear)
	}

	if req.IgnorePaging {
		data.Meta = pagination.CreateMeta(totalData, totalData, 1)
	} else {
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	}

	return
}

func (r *AcademicYearRepositoryPostgreSQL) GetAll(ctx context.Context, req model.StandardRequest) (data []AcademicYearDTO, err error) {
	var searchParams []interface{}
	var searchBuff bytes.Buffer
	searchBuff.WriteString(" WHERE coalesce(ay.is_deleted, false) = false")
	searchBuff.WriteString(" ORDER BY ay.name ASC ")
	query := r.DB.Read.Rebind(AcademicYearQuery.SelectDTO + searchBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, searchParams...)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Academic Year Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var academicYear AcademicYearDTO
		err = rows.StructScan(&academicYear)

		if err != nil {
			return
		}
		data = append(data, academicYear)
	}
	return
}

func (r *AcademicYearRepositoryPostgreSQL) ResolveByID(ctx context.Context, id int) (AcademicYear AcademicYear, err error) {
	err = r.DB.Read.GetContext(ctx, &AcademicYear, AcademicYearQuery.Select+" WHERE id=$1", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *AcademicYearRepositoryPostgreSQL) ResolveByCode(ctx context.Context, code string, institutionID int) (AcademicYear AcademicYear, err error) {
	err = r.DB.Read.GetContext(ctx, &AcademicYear, AcademicYearQuery.Select+" WHERE code=$1 AND institution_id=$2 AND coalesce(is_deleted, false)=false", code, institutionID)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (r *AcademicYearRepositoryPostgreSQL) DeleteByID(ctx context.Context, id int) (err error) {
	_, err = r.DB.Write.QueryContext(ctx, AcademicYearQuery.Delete, id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *AcademicYearRepositoryPostgreSQL) ExistByField(ctx context.Context, params model.ExistByFieldParams) (bool, error) {
	var idResult int
	where := fmt.Sprintf(" WHERE UPPER(%s) = UPPER($1) AND institution_id = $2 AND coalesce(is_deleted, false) = false", params.Field)
	if params.ExcludeID != 0 {
		where += fmt.Sprintf(" AND id <> %d", params.ExcludeID)
	}

	err := r.DB.Read.GetContext(ctx, &idResult, AcademicYearQuery.Exist+where, params.Value, params.InstitutionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.ErrorWithStack(err)
		return false, err
	}
	return true, nil
}
