package internship

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"lms-be/infras"
	"lms-be/shared/logger"
)

var TaskFileQuery = struct {
	Insert string
}{
	Insert: `INSERT INTO task_files (id, task_id, file_url, uploaded_by, created_at)
		VALUES (:id, :task_id, :file_url, :uploaded_by, :created_at) RETURNING id`,
}

type TaskFileRepository interface {
	Create(ctx context.Context, data *TaskFile) error
	CreateTx(ctx context.Context, tx *sqlx.Tx, data *TaskFile) error
}

type TaskFileRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideTaskFileRepositoryPostgreSQL(db *infras.PostgresqlConn) *TaskFileRepositoryPostgreSQL {
	s := new(TaskFileRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *TaskFileRepositoryPostgreSQL) Create(ctx context.Context, data *TaskFile) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, TaskFileQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	defer stmt.Close()

	if err = stmt.QueryRowContext(ctx, data).Scan(&data.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

func (r *TaskFileRepositoryPostgreSQL) CreateTx(ctx context.Context, tx *sqlx.Tx, data *TaskFile) error {
	stmt, err := tx.PrepareNamedContext(ctx, TaskFileQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	defer stmt.Close()

	if err = stmt.QueryRowContext(ctx, data).Scan(&data.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}
