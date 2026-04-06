package internship

import (
	"bytes"
	"context"
	"database/sql"

	"github.com/gofrs/uuid/v5"

	"lms-be/infras"
	"lms-be/shared/logger"
)

var RegistrationQuery = struct {
	Select        string
	SelectDTO     string
	Insert        string
	ExistByUserID string
}{
	Select:    `SELECT id, user_id, university, major, cv_file_path, status, created_at FROM internship_registrations `,
	SelectDTO: `SELECT id, user_id, university, major, cv_file_path, status, created_at FROM internship_registrations `,
	Insert: `INSERT INTO internship_registrations
		(id, user_id, university, major, cv_file_path, status, created_at)
		VALUES
		(:id, :user_id, :university, :major, :cv_file_path, :status, :created_at) RETURNING id`,
	ExistByUserID: `SELECT id FROM internship_registrations`,
}

type RegistrationRepository interface {
	Create(ctx context.Context, data *Registration) error
	GetAll(ctx context.Context) (data []RegistrationDTO, err error)
	ExistByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
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

func (r *RegistrationRepositoryPostgreSQL) GetAll(ctx context.Context) (data []RegistrationDTO, err error) {
	var searchBuff bytes.Buffer
	searchBuff.WriteString(" ORDER BY created_at DESC ")

	query := r.DB.Read.Rebind(RegistrationQuery.SelectDTO + searchBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query)
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
