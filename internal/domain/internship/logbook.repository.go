package internship

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"lms-be/infras"
	"lms-be/shared/logger"
	"lms-be/shared/pagination"
)

var LogbookQuery = struct {
	Select           string
	SelectStudentDTO string
	SelectMentorDTO  string
	Insert           string
	UpdateStatus     string
	Update           string
	Count            string
}{
	Select: `SELECT id, student_id, log_date, activities, blockers, plan_tomorrow, evidence_url, status, notes, submitted_at, reviewed_at, reviewed_by FROM logbooks `,
	SelectStudentDTO: `SELECT l.id, l.student_id, s.name AS student_name, ma.mentor_id, mentor.name AS mentor_name, l.log_date, l.activities, l.blockers, l.plan_tomorrow, l.evidence_url, l.status, l.notes, l.submitted_at, l.reviewed_at, l.reviewed_by
		FROM logbooks l
		LEFT JOIN auth_user s ON s.id = l.student_id
		LEFT JOIN mentor_assignments ma ON ma.student_id = l.student_id AND ma.is_active = true
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id`,
	SelectMentorDTO: `SELECT l.id, l.student_id, s.name AS student_name, ma.mentor_id, mentor.name AS mentor_name, l.log_date, l.activities, l.blockers, l.plan_tomorrow, l.evidence_url, l.status, l.notes, l.submitted_at, l.reviewed_at, l.reviewed_by
		FROM logbooks l
		INNER JOIN mentor_assignments ma ON ma.student_id = l.student_id AND ma.is_active = true
		LEFT JOIN auth_user s ON s.id = l.student_id
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id`,
	Insert: `INSERT INTO logbooks
		(id, student_id, log_date, activities, blockers, plan_tomorrow, evidence_url, status, submitted_at)
		VALUES
		(:id, :student_id, :log_date, :activities, :blockers, :plan_tomorrow, :evidence_url, :status, :submitted_at) RETURNING id`,
	UpdateStatus: `UPDATE logbooks SET status = :status, notes = :notes, reviewed_at = :reviewed_at, reviewed_by = :reviewed_by WHERE id = :id`,
	Update: `UPDATE logbooks SET log_date = :log_date, activities = :activities, blockers = :blockers, plan_tomorrow = :plan_tomorrow, evidence_url = :evidence_url WHERE id = :id`,
	Count:  `SELECT count(*) FROM logbooks `,
}

type LogbookRepository interface {
	Create(ctx context.Context, data *Logbook) error
	GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []LogbookDTO, err error)
	GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []LogbookDTO, err error)
	ResolveByID(ctx context.Context, id uuid.UUID) (data Logbook, err error)
	UpdateStatus(ctx context.Context, data Logbook) error
	Update(ctx context.Context, data Logbook) error
	ResolveAllByStudentID(ctx context.Context, studentID uuid.UUID, req RequestLogbookListFormat) (data pagination.Response, err error)
	ResolveAllByMentorID(ctx context.Context, mentorID uuid.UUID, req RequestLogbookListFormat) (data pagination.Response, err error)
}

type LogbookRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideLogbookRepositoryPostgreSQL(db *infras.PostgresqlConn) *LogbookRepositoryPostgreSQL {
	s := new(LogbookRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *LogbookRepositoryPostgreSQL) Create(ctx context.Context, data *Logbook) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, LogbookQuery.Insert)
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

func (r *LogbookRepositoryPostgreSQL) GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []LogbookDTO, err error) {
	query := r.DB.Read.Rebind(LogbookQuery.SelectStudentDTO + " WHERE l.student_id = ? ORDER BY l.log_date DESC, l.submitted_at DESC")
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
		var logbook LogbookDTO
		err = rows.StructScan(&logbook)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data = append(data, logbook)
	}
	return
}

func (r *LogbookRepositoryPostgreSQL) GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []LogbookDTO, err error) {
	query := r.DB.Read.Rebind(LogbookQuery.SelectMentorDTO + " WHERE ma.mentor_id = ? ORDER BY l.log_date DESC, l.submitted_at DESC")
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
		var logbook LogbookDTO
		err = rows.StructScan(&logbook)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data = append(data, logbook)
	}
	return
}

func (r *LogbookRepositoryPostgreSQL) ResolveByID(ctx context.Context, id uuid.UUID) (data Logbook, err error) {
	err = r.DB.Read.GetContext(ctx, &data, r.DB.Read.Rebind(LogbookQuery.Select+" WHERE id = ?"), id)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (r *LogbookRepositoryPostgreSQL) UpdateStatus(ctx context.Context, data Logbook) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, LogbookQuery.UpdateStatus)
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

func (r *LogbookRepositoryPostgreSQL) Update(ctx context.Context, data Logbook) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, LogbookQuery.Update)
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

func (r *LogbookRepositoryPostgreSQL) ResolveAllByStudentID(ctx context.Context, studentID uuid.UUID, req RequestLogbookListFormat) (data pagination.Response, err error) {
	var totalData int
	var searchBuff bytes.Buffer
	searchBuff.WriteString(" WHERE l.student_id = ? ")

	if req.Search != "" {
		searchBuff.WriteString(fmt.Sprintf(" AND (l.activities ILIKE '%%%s%%' OR l.blockers ILIKE '%%%s%%' OR l.plan_tomorrow ILIKE '%%%s%%') ", req.Search, req.Search, req.Search))
	}

	if req.Status != "" {
		searchBuff.WriteString(fmt.Sprintf(" AND l.status = '%s' ", req.Status))
	}

	countQuery := r.DB.Read.Rebind(LogbookQuery.Count + " l " + searchBuff.String())
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

	searchBuff.WriteString(" ORDER BY l.log_date DESC, l.submitted_at DESC ")
	searchBuff.WriteString(" LIMIT ? OFFSET ? ")

	offset := (req.PageNumber - 1) * req.PageSize
	query := r.DB.Read.Rebind(LogbookQuery.SelectStudentDTO + searchBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, studentID, req.PageSize, offset)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var logbook LogbookDTO
		err = rows.StructScan(&logbook)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data.Items = append(data.Items, logbook)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *LogbookRepositoryPostgreSQL) ResolveAllByMentorID(ctx context.Context, mentorID uuid.UUID, req RequestLogbookListFormat) (data pagination.Response, err error) {
	var totalData int
	var filterBuff bytes.Buffer
	filterBuff.WriteString(" WHERE ma.mentor_id = ? ")

	if req.Search != "" {
		filterBuff.WriteString(fmt.Sprintf(" AND (l.activities ILIKE '%%%s%%' OR l.blockers ILIKE '%%%s%%' OR l.plan_tomorrow ILIKE '%%%s%%') ", req.Search, req.Search, req.Search))
	}

	if req.Status != "" {
		filterBuff.WriteString(fmt.Sprintf(" AND l.status = '%s' ", req.Status))
	}

	joinSQL := " INNER JOIN mentor_assignments ma ON ma.student_id = l.student_id AND ma.is_active = true "
	countQuery := r.DB.Read.Rebind(LogbookQuery.Count + " l " + joinSQL + filterBuff.String())
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

	filterBuff.WriteString(" ORDER BY l.log_date DESC, l.submitted_at DESC ")
	filterBuff.WriteString(" LIMIT ? OFFSET ? ")

	offset := (req.PageNumber - 1) * req.PageSize
	query := r.DB.Read.Rebind(LogbookQuery.SelectMentorDTO + filterBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, mentorID, req.PageSize, offset)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var logbook LogbookDTO
		err = rows.StructScan(&logbook)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data.Items = append(data.Items, logbook)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}
