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
	PersonnelQuery = struct {
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
		Select: `SELECT id, name, personnel_type, department_id, company_id, photo_face_template, is_active, created_at, created_by, updated_at, updated_by, deleted_at, is_deleted FROM m_personnel `,
		SelectDTO: `SELECT p.id, p.name, p.personnel_type, p.department_id, p.company_id, p.photo_face_template, p.is_active, c.name as company_name, d.name as department_name
			FROM m_personnel p
			LEFT JOIN m_company c ON p.company_id = c.id
			LEFT JOIN m_department d ON p.department_id = d.id `,
		Insert: `INSERT INTO m_personnel
			(name, personnel_type, department_id, company_id, photo_face_template, is_active, created_at, created_by)
			VALUES
			(:name, :personnel_type, :department_id, :company_id, :photo_face_template, :is_active, :created_at, :created_by)
			RETURNING id`,
		Update:                 `UPDATE m_personnel SET name = :name, personnel_type = :personnel_type, department_id = :department_id, company_id = :company_id, photo_face_template = :photo_face_template, is_active = :is_active, updated_at = :updated_at, updated_by = :updated_by, deleted_at = :deleted_at, is_deleted = :is_deleted WHERE id = :id`,
		Delete:                 `DELETE FROM m_personnel WHERE id = $1`,
		Count:                  `SELECT count(p.id) FROM m_personnel p `,
		Exist:                  `SELECT p.id FROM m_personnel p`,
		InsertBatch:            `INSERT INTO m_personnel (name, personnel_type, department_id, company_id, photo_face_template, is_active, created_at, created_by) VALUES `,
		InsertBatchPlaceholder: `(:name, :personnel_type, :department_id, :company_id, :photo_face_template, :is_active, :created_at, :created_by)`,
		UpsertBatchOnConflict:  `DO UPDATE SET name = EXCLUDED.name, personnel_type = EXCLUDED.personnel_type, department_id = EXCLUDED.department_id, company_id = EXCLUDED.company_id, photo_face_template = EXCLUDED.photo_face_template, is_active = EXCLUDED.is_active, updated_at = EXCLUDED.created_at, updated_by = EXCLUDED.created_by`,
	}
)

type PersonnelRepository interface {
	Create(ctx context.Context, data *Personnel) error
	CreateBatchTx(ctx context.Context, tx *sqlx.Tx, data []Personnel) error
	UpsertBatchTx(ctx context.Context, tx *sqlx.Tx, data []Personnel) error
	Update(ctx context.Context, data Personnel) error
	GetAll(ctx context.Context, req model.StandardRequest) (data []PersonnelDTO, err error)
	ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(ctx context.Context, id int) (Personnel Personnel, err error)
	DeleteByID(ctx context.Context, id int) (err error)
	ExistByField(ctx context.Context, params model.ExistByFieldParams) (bool, error)
	ExistDepartment(ctx context.Context, departmentID int) (bool, error)
	ExistCompany(ctx context.Context, companyID int) (bool, error)
	GetDepartmentNameByID(ctx context.Context, departmentID int) (string, error)
	GetCompanyNameByID(ctx context.Context, companyID int) (string, error)
}

type PersonnelRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvidePersonnelRepositoryPostgreSQL(db *infras.PostgresqlConn) *PersonnelRepositoryPostgreSQL {
	s := new(PersonnelRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *PersonnelRepositoryPostgreSQL) Create(ctx context.Context, data *Personnel) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, PersonnelQuery.Insert)
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

func (r *PersonnelRepositoryPostgreSQL) CreateBatchTx(ctx context.Context, tx *sqlx.Tx, data []Personnel) error {
	return r.ProcessBatchInsertTx(ctx, tx, data, false)
}

func (r *PersonnelRepositoryPostgreSQL) UpsertBatchTx(ctx context.Context, tx *sqlx.Tx, data []Personnel) error {
	return r.ProcessBatchInsertTx(ctx, tx, data, true)
}

func (r *PersonnelRepositoryPostgreSQL) ProcessBatchInsertTx(ctx context.Context, tx *sqlx.Tx, data []Personnel, upsert bool) error {
	if len(data) == 0 {
		return nil
	}

	suffix := " ON CONFLICT (name) DO NOTHING "
	if upsert {
		suffix = " ON CONFLICT (name) " + PersonnelQuery.UpsertBatchOnConflict
	}

	config := repository.BatchInsertConfig{
		BaseQuery:       PersonnelQuery.InsertBatch,
		Placeholder:     PersonnelQuery.InsertBatchPlaceholder,
		ReturningSuffix: suffix,
		MaxBatchSize:    1000,
	}

	mapFunc := func(item Personnel) map[string]interface{} {
		return map[string]interface{}{
			"name":                item.Name,
			"personnel_type":      item.PersonnelType,
			"department_id":       item.DepartmentID,
			"company_id":          item.CompanyID,
			"photo_face_template": item.PhotoFaceTemplate,
			"is_active":           item.IsActive,
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
		_, err := tx.ExecContext(ctx, chunk.Query, chunk.Args...)
		if err != nil {
			logger.ErrorWithStack(err)
			return err
		}
	}

	return nil
}

func (p *PersonnelRepositoryPostgreSQL) Update(ctx context.Context, data Personnel) error {
	stmt, err := p.DB.Write.PrepareNamedContext(ctx, PersonnelQuery.Update)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, data)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return err
}

func (r *PersonnelRepositoryPostgreSQL) ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(p.is_deleted, false) = false")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" p.name, c.name, d.name ILIKE ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind("SELECT COUNT(x.id) FROM (" + PersonnelQuery.SelectDTO + searchRoleBuff.String() + ")x")

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

	searchRoleBuff.WriteString(" ORDER BY " + ColumnMapPersonnel[req.SortBy].(string) + " " + req.SortType)

	if !req.IgnorePaging {
		offset := (req.PageNumber - 1) * req.PageSize
		searchRoleBuff.WriteString(" LIMIT ? OFFSET ? ")
		searchParams = append(searchParams, req.PageSize)
		searchParams = append(searchParams, offset)
	}

	searchPersonnelQuery := searchRoleBuff.String()
	searchPersonnelQuery = r.DB.Read.Rebind(PersonnelQuery.SelectDTO + searchPersonnelQuery)

	rows, err := r.DB.Read.QueryxContext(ctx, searchPersonnelQuery, searchParams...)

	if err != nil {
		return
	}

	for rows.Next() {
		var personnelDTO PersonnelDTO
		err = rows.StructScan(&personnelDTO)
		if err != nil {
			return
		}
		data.Items = append(data.Items, personnelDTO)
	}

	if req.IgnorePaging {
		data.Meta = pagination.CreateMeta(totalData, totalData, 1)
	} else {
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	}

	return
}

func (r *PersonnelRepositoryPostgreSQL) GetAll(ctx context.Context, req model.StandardRequest) (data []PersonnelDTO, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(p.is_deleted, false) = false")
	searchRoleBuff.WriteString(" ORDER BY p.created_at DESC ")
	query := r.DB.Read.Rebind(PersonnelQuery.SelectDTO + searchRoleBuff.String())
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
		var personnelDTO PersonnelDTO
		err = rows.StructScan(&personnelDTO)

		if err != nil {
			return
		}
		data = append(data, personnelDTO)
	}

	return
}

func (r *PersonnelRepositoryPostgreSQL) ResolveByID(ctx context.Context, id int) (personnel Personnel, err error) {
	err = r.DB.Read.GetContext(ctx, &personnel, PersonnelQuery.Select+` WHERE id = $1`, id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *PersonnelRepositoryPostgreSQL) DeleteByID(ctx context.Context, id int) (err error) {
	_, err = r.DB.Write.ExecContext(ctx, PersonnelQuery.Delete, id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *PersonnelRepositoryPostgreSQL) ExistByField(ctx context.Context, params model.ExistByFieldParams) (bool, error) {
	var idResult int
	where := fmt.Sprintf(" WHERE UPPER(%s) = UPPER($1) AND coalesce(is_deleted, false) = false", params.Field)
	if params.ExcludeID != 0 {
		where += fmt.Sprintf(" AND id <> %d", params.ExcludeID)
	}

	err := r.DB.Read.GetContext(ctx, &idResult, PersonnelQuery.Exist+where, params.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.ErrorWithStack(err)
		return false, err
	}

	return true, nil
}

func (r *PersonnelRepositoryPostgreSQL) ExistDepartment(ctx context.Context, departmentID int) (bool, error) {
	query := `SELECT id FROM m_department WHERE id = $1 AND is_deleted = false`
	row := r.DB.Read.QueryRowContext(ctx, query, departmentID)
	var id int
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		logger.ErrorWithStack(err)
		return false, err
	}
	return err != sql.ErrNoRows, nil
}

func (r *PersonnelRepositoryPostgreSQL) ExistCompany(ctx context.Context, companyID int) (bool, error) {
	query := `SELECT id FROM m_company WHERE id = $1 AND is_deleted = false`
	row := r.DB.Read.QueryRowContext(ctx, query, companyID)
	var id int
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		logger.ErrorWithStack(err)
		return false, err
	}
	return err != sql.ErrNoRows, nil
}

func (r *PersonnelRepositoryPostgreSQL) GetDepartmentNameByID(ctx context.Context, departmentID int) (string, error) {
	query := `SELECT name FROM m_department WHERE id = $1 AND is_deleted = false`
	row := r.DB.Read.QueryRowContext(ctx, query, departmentID)
	var name string
	err := row.Scan(&name)
	if err != nil {
		logger.ErrorWithStack(err)
		return "", err
	}
	return name, nil
}

func (r *PersonnelRepositoryPostgreSQL) GetCompanyNameByID(ctx context.Context, companyID int) (string, error) {
	query := `SELECT name FROM m_company WHERE id = $1 AND is_deleted = false`
	row := r.DB.Read.QueryRowContext(ctx, query, companyID)
	var name string
	err := row.Scan(&name)
	if err != nil {
		logger.ErrorWithStack(err)
		return "", err
	}
	return name, nil
}
