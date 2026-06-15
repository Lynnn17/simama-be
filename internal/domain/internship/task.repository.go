package internship

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"

	"lms-be/infras"
	"lms-be/shared/logger"
	"lms-be/shared/pagination"
)

var TaskQuery = struct {
	Select             string
	SelectDTO          string
	Insert             string
	Update             string
	UpdateSubmitted    string
	UpdateGrade        string
	ExistByStudentID   string
	ExistByMentorID    string
	ResolveByID        string
	SelectFileByTaskID string
	Count              string
}{
	Select: `SELECT id, mentor_id, student_id, title, description, deadline, status, grade, feedback, submission_url, created_at, updated_at FROM tasks `,
	SelectDTO: `SELECT t.id, t.mentor_id, mentor.name AS mentor_name, t.student_id, student.name AS student_name, t.title, t.description, t.deadline, t.status, t.grade, t.feedback, t.submission_url, t.created_at, t.updated_at, tf.id AS latest_file_id, tf.file_url AS latest_file_url
		FROM tasks t
		LEFT JOIN auth_user mentor ON mentor.id = t.mentor_id
		LEFT JOIN auth_user student ON student.id = t.student_id
		LEFT JOIN LATERAL (
			SELECT tf1.id, tf1.file_url
			FROM task_files tf1
			WHERE tf1.task_id = t.id
			ORDER BY tf1.created_at DESC
			LIMIT 1
		) tf ON TRUE`,
	Insert: `INSERT INTO tasks (id, mentor_id, student_id, title, description, deadline, status, submission_url, created_at, updated_at)
		VALUES (:id, :mentor_id, :student_id, :title, :description, :deadline, :status, :submission_url, :created_at, :updated_at) RETURNING id`,
	Update:             `UPDATE tasks SET mentor_id = :mentor_id, student_id = :student_id, title = :title, description = :description, deadline = :deadline, status = :status, grade = :grade, feedback = :feedback, submission_url = :submission_url, updated_at = :updated_at WHERE id = :id`,
	UpdateSubmitted:    `UPDATE tasks SET status = :status, submission_url = :submission_url, updated_at = :updated_at WHERE id = :id`,
	UpdateGrade:        `UPDATE tasks SET status = :status, grade = :grade, feedback = :feedback, updated_at = :updated_at WHERE id = :id`,
	ExistByStudentID:   `SELECT id FROM tasks`,
	ExistByMentorID:    `SELECT id FROM tasks`,
	ResolveByID:        `SELECT id, mentor_id, student_id, title, description, deadline, status, grade, feedback, submission_url, created_at, updated_at FROM tasks`,
	SelectFileByTaskID: `SELECT id, task_id, file_url, uploaded_by, created_at FROM task_files `,
	Count:              `SELECT count(*) FROM tasks `,
}

type TaskRepository interface {
	Create(ctx context.Context, data *Task) error
	CreateTx(ctx context.Context, tx *sqlx.Tx, data *Task) error
	GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []TaskDTO, err error)
	GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []TaskDTO, err error)
	ResolveAllByStudentID(ctx context.Context, studentID uuid.UUID, req RequestTaskListFormat) (data pagination.Response, err error)
	ResolveAllByMentorID(ctx context.Context, mentorID uuid.UUID, req RequestTaskListFormat) (data pagination.Response, err error)
	ResolveByID(ctx context.Context, id uuid.UUID) (data Task, err error)
	Update(ctx context.Context, data Task) error
	UpdateSubmittedTx(ctx context.Context, tx *sqlx.Tx, data Task) error
	UpdateGradeTx(ctx context.Context, tx *sqlx.Tx, data Task) error
	UpdateSubmitted(ctx context.Context, data Task) error
	UpdateGrade(ctx context.Context, data Task) error
	ExistByStudentID(ctx context.Context, studentID uuid.UUID) (bool, error)
	ExistByMentorID(ctx context.Context, mentorID uuid.UUID) (bool, error)
}

type TaskRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideTaskRepositoryPostgreSQL(db *infras.PostgresqlConn) *TaskRepositoryPostgreSQL {
	s := new(TaskRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *TaskRepositoryPostgreSQL) Create(ctx context.Context, data *Task) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, TaskQuery.Insert)
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

func (r *TaskRepositoryPostgreSQL) CreateTx(ctx context.Context, tx *sqlx.Tx, data *Task) error {
	stmt, err := tx.PrepareNamedContext(ctx, TaskQuery.Insert)
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

func (r *TaskRepositoryPostgreSQL) GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []TaskDTO, err error) {
	query := r.DB.Read.Rebind(TaskQuery.SelectDTO + " WHERE t.student_id = ? ORDER BY t.created_at DESC")
	rows, err := r.DB.Read.QueryxContext(ctx, query, studentID)
	if err == sql.ErrNoRows {
		return data, nil
	}
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task TaskDTO
		err = rows.StructScan(&task)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data = append(data, task)
	}
	return
}

func (r *TaskRepositoryPostgreSQL) GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []TaskDTO, err error) {
	query := r.DB.Read.Rebind(TaskQuery.SelectDTO + " WHERE t.mentor_id = ? ORDER BY t.created_at DESC")
	rows, err := r.DB.Read.QueryxContext(ctx, query, mentorID)
	if err == sql.ErrNoRows {
		return data, nil
	}
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task TaskDTO
		err = rows.StructScan(&task)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data = append(data, task)
	}
	return
}

func (r *TaskRepositoryPostgreSQL) ResolveByID(ctx context.Context, id uuid.UUID) (data Task, err error) {
	err = r.DB.Read.GetContext(ctx, &data, r.DB.Read.Rebind(TaskQuery.ResolveByID+" WHERE id = ?"), id)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (r *TaskRepositoryPostgreSQL) Update(ctx context.Context, data Task) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, TaskQuery.Update)
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

func (r *TaskRepositoryPostgreSQL) UpdateSubmitted(ctx context.Context, data Task) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, TaskQuery.UpdateSubmitted)
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

func (r *TaskRepositoryPostgreSQL) UpdateSubmittedTx(ctx context.Context, tx *sqlx.Tx, data Task) error {
	stmt, err := tx.PrepareNamedContext(ctx, TaskQuery.UpdateSubmitted)
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

func (r *TaskRepositoryPostgreSQL) UpdateGrade(ctx context.Context, data Task) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, TaskQuery.UpdateGrade)
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

func (r *TaskRepositoryPostgreSQL) UpdateGradeTx(ctx context.Context, tx *sqlx.Tx, data Task) error {
	stmt, err := tx.PrepareNamedContext(ctx, TaskQuery.UpdateGrade)
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

func (r *TaskRepositoryPostgreSQL) ExistByStudentID(ctx context.Context, studentID uuid.UUID) (bool, error) {
	var idResult uuid.UUID
	err := r.DB.Read.GetContext(ctx, &idResult, r.DB.Read.Rebind(TaskQuery.ExistByStudentID+" WHERE student_id = ?"), studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.ErrorWithStack(err)
		return false, err
	}
	return true, nil
}

func (r *TaskRepositoryPostgreSQL) ExistByMentorID(ctx context.Context, mentorID uuid.UUID) (bool, error) {
	var idResult uuid.UUID
	err := r.DB.Read.GetContext(ctx, &idResult, r.DB.Read.Rebind(TaskQuery.ExistByMentorID+" WHERE mentor_id = ?"), mentorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.ErrorWithStack(err)
		return false, err
	}
	return true, nil
}

func (r *TaskRepositoryPostgreSQL) ResolveAllByStudentID(ctx context.Context, studentID uuid.UUID, req RequestTaskListFormat) (data pagination.Response, err error) {
	var totalData int
	var filterBuff bytes.Buffer
	filterBuff.WriteString(" WHERE t.student_id = ? ")

	if req.Search != "" {
		filterBuff.WriteString(fmt.Sprintf(" AND (t.title ILIKE '%%%s%%' OR t.description ILIKE '%%%s%%') ", req.Search, req.Search))
	}

	if req.Status != "" {
		filterBuff.WriteString(fmt.Sprintf(" AND t.status = '%s' ", req.Status))
	}

	if req.Date != "" {
		filterBuff.WriteString(fmt.Sprintf(" AND DATE(t.deadline) = '%s' ", req.Date))
	}

	countQuery := r.DB.Read.Rebind(TaskQuery.Count + " t " + filterBuff.String())
	err = r.DB.Read.QueryRowContext(ctx, countQuery, studentID).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		data.Items = make([]interface{}, 0)
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
		return
	}

	filterBuff.WriteString(" ORDER BY t.created_at DESC ")
	filterBuff.WriteString(" LIMIT ? OFFSET ? ")

	offset := (req.PageNumber - 1) * req.PageSize
	query := r.DB.Read.Rebind(TaskQuery.SelectDTO + filterBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, studentID, req.PageSize, offset)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task TaskDTO
		err = rows.StructScan(&task)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data.Items = append(data.Items, task)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *TaskRepositoryPostgreSQL) ResolveAllByMentorID(ctx context.Context, mentorID uuid.UUID, req RequestTaskListFormat) (data pagination.Response, err error) {
	var totalData int
	var filterBuff bytes.Buffer
	filterBuff.WriteString(" WHERE t.mentor_id = ? ")

	if req.Search != "" {
		filterBuff.WriteString(fmt.Sprintf(" AND (t.title ILIKE '%%%s%%') ", req.Search))
	}

	if req.Status != "" {
		filterBuff.WriteString(fmt.Sprintf(" AND t.status = '%s' ", req.Status))
	}

	if req.Date != "" {
		filterBuff.WriteString(fmt.Sprintf(" AND DATE(t.deadline) = '%s' ", req.Date))
	}

	joinSQL := ""
	if req.StudentSearch != "" {
		joinSQL = " LEFT JOIN auth_user student ON student.id = t.student_id "
		filterBuff.WriteString(fmt.Sprintf(" AND student.name ILIKE '%%%s%%' ", req.StudentSearch))
	}

	countQuery := r.DB.Read.Rebind(TaskQuery.Count + " t " + joinSQL + filterBuff.String())
	err = r.DB.Read.QueryRowContext(ctx, countQuery, mentorID).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		data.Items = make([]interface{}, 0)
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
		return
	}

	filterBuff.WriteString(" ORDER BY t.created_at DESC ")
	filterBuff.WriteString(" LIMIT ? OFFSET ? ")

	offset := (req.PageNumber - 1) * req.PageSize
	query := r.DB.Read.Rebind(TaskQuery.SelectDTO + filterBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, mentorID, req.PageSize, offset)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task TaskDTO
		err = rows.StructScan(&task)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data.Items = append(data.Items, task)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}
