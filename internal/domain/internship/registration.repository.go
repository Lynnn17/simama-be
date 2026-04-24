package internship

import (
	"bytes"
	"context"
	"database/sql"

	"github.com/gofrs/uuid/v5"

	"lms-be/infras"
	"lms-be/shared/logger"
	"lms-be/shared/pagination"
)

var RegistrationQuery = struct {
	Select        string
	SelectDTO     string
	Insert        string
	ExistByUserID string
	Count         string
	UpdateStatus  string
	UpdateUserID  string
	ResolveByID   string
}{
	Select:    `SELECT id, user_id, full_name, university, major, semester, phone, email, period, cv_file_path, status, applied_at, reviewed_at, reviewed_by FROM internship_registrations `,
	SelectDTO: `SELECT id, user_id, full_name, university, major, semester, phone, email, period, cv_file_path, status, applied_at, reviewed_at, reviewed_by FROM internship_registrations `,
	Insert: `INSERT INTO internship_registrations 
		(id, user_id, full_name, university, major, semester, phone, email, period, cv_file_path, status, applied_at, reviewed_at, reviewed_by)
		VALUES
		(:id, :user_id, :full_name, :university, :major, :semester, :phone, :email, :period, :cv_file_path, :status, :applied_at, :reviewed_at, :reviewed_by) RETURNING id`,
	ExistByUserID: `SELECT id FROM internship_registrations`,
	Count:         `SELECT count(id) FROM internship_registrations `,
	UpdateStatus:  `UPDATE internship_registrations SET status = :status, reviewed_at = :reviewed_at, reviewed_by = :reviewed_by WHERE id = :id`,
	UpdateUserID:  `UPDATE internship_registrations SET user_id = :user_id WHERE id = :id`,
	ResolveByID:   `SELECT id, user_id, full_name, university, major, semester, phone, email, period, cv_file_path, status, applied_at, reviewed_at, reviewed_by FROM internship_registrations`,
}

type RegistrationRepository interface {
	Create(ctx context.Context, data *Registration) error
	GetAll(ctx context.Context, req RequestRegistrationListFormat) (data []RegistrationDTO, err error)
	ResolveAll(ctx context.Context, req RequestRegistrationListFormat) (data pagination.Response, err error)
	ResolveByID(ctx context.Context, id uuid.UUID) (data Registration, err error)
	ExistByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
	UpdateStatus(ctx context.Context, data Registration) error
	UpdateUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type RegistrationRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideRegistrationRepositoryPostgreSQL(db *infras.PostgresqlConn) *RegistrationRepositoryPostgreSQL {
	s := new(RegistrationRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *RegistrationRepositoryPostgreSQL) Create(ctx context.Context, data *Registration) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, RegistrationQuery.Insert)
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

func (r *RegistrationRepositoryPostgreSQL) GetAll(ctx context.Context, req RequestRegistrationListFormat) (data []RegistrationDTO, err error) {
	var searchParams []interface{}
	var searchBuff bytes.Buffer
	searchBuff.WriteString(" WHERE 1=1 ")

	if req.Status != "" {
		searchBuff.WriteString(" AND status = ? ")
		searchParams = append(searchParams, req.Status)
	}

	searchBuff.WriteString(" ORDER BY applied_at DESC ")

	query := r.DB.Read.Rebind(RegistrationQuery.SelectDTO + searchBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, searchParams...)
	if err == sql.ErrNoRows {
		return data, nil
	}
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var registration RegistrationDTO
		err = rows.StructScan(&registration)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data = append(data, registration)
	}
	return
}

func (r *RegistrationRepositoryPostgreSQL) UpdateStatus(ctx context.Context, data Registration) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, RegistrationQuery.UpdateStatus)
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

func (r *RegistrationRepositoryPostgreSQL) UpdateUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := r.DB.Write.Rebind("UPDATE internship_registrations SET user_id = ? WHERE id = ?")
	_, err := r.DB.Write.ExecContext(ctx, query, userID, id)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return err
}

func (r *RegistrationRepositoryPostgreSQL) ResolveByID(ctx context.Context, id uuid.UUID) (data Registration, err error) {
	err = r.DB.Read.GetContext(ctx, &data, r.DB.Read.Rebind(RegistrationQuery.ResolveByID+" WHERE id = ?"), id)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (r *RegistrationRepositoryPostgreSQL) ResolveAll(ctx context.Context, req RequestRegistrationListFormat) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchBuff bytes.Buffer
	searchBuff.WriteString(" WHERE 1=1 ")

	if req.Status != "" {
		searchBuff.WriteString(" AND status = ? ")
		searchParams = append(searchParams, req.Status)
	}

	var totalData int
	countQuery := r.DB.Read.Rebind(RegistrationQuery.Count + searchBuff.String())
	err = r.DB.Read.QueryRowContext(ctx, countQuery, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		data.Items = make([]interface{}, 0)
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
		return
	}

	searchBuff.WriteString(" ORDER BY applied_at DESC ")
	searchBuff.WriteString(" LIMIT ? OFFSET ? ")

	offset := (req.PageNumber - 1) * req.PageSize
	query := r.DB.Read.Rebind(RegistrationQuery.SelectDTO + searchBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, append(searchParams, req.PageSize, offset)...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var registration RegistrationDTO
		err = rows.StructScan(&registration)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data.Items = append(data.Items, registration)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *RegistrationRepositoryPostgreSQL) ExistByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var idResult uuid.UUID
	query := r.DB.Read.Rebind(RegistrationQuery.ExistByUserID + " WHERE user_id = ?")
	err := r.DB.Read.GetContext(ctx, &idResult, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.ErrorWithStack(err)
		return false, err
	}
	return true, nil
}
