package internship

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"

	"lms-be/infras"
	"lms-be/shared/logger"
)

var TaskFileQuery = struct {
	Insert         string
	SelectByTaskID string
}{
	Insert:         `INSERT INTO task_files (id, task_id, file_url, uploaded_by, created_at) VALUES (:id, :task_id, :file_url, :uploaded_by, :created_at) RETURNING id`,
	SelectByTaskID: `SELECT id, task_id, file_url, uploaded_by, created_at FROM task_files WHERE task_id = ? ORDER BY created_at DESC`,
}

type TaskFileRepository interface {
	Create(ctx context.Context, data *TaskFile) error
	CreateTx(ctx context.Context, tx *sqlx.Tx, data *TaskFile) error
	ResolveByTaskID(ctx context.Context, taskID uuid.UUID) (data []TaskFile, err error)
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

func (r *TaskFileRepositoryPostgreSQL) ResolveByTaskID(ctx context.Context, taskID uuid.UUID) (data []TaskFile, err error) {
	err = r.DB.Read.SelectContext(ctx, &data, r.DB.Read.Rebind(TaskFileQuery.SelectByTaskID), taskID)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}
