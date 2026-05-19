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
	Update           string
	Count            string
}{
	Select: `SELECT id, student_id, log_date, activities, blockers, plan_tomorrow, evidence_url, progress_status, submitted_at FROM logbooks `,
	SelectStudentDTO: `SELECT l.id, l.student_id, s.name AS student_name, ma.mentor_id, mentor.name AS mentor_name, l.log_date, l.activities, l.blockers, l.plan_tomorrow, l.evidence_url, l.progress_status, l.submitted_at
		FROM logbooks l
		LEFT JOIN auth_user s ON s.id = l.student_id
		LEFT JOIN mentor_assignments ma ON ma.student_id = l.student_id AND ma.is_active = true
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id`,
	SelectMentorDTO: `SELECT l.id, l.student_id, s.name AS student_name, ma.mentor_id, mentor.name AS mentor_name, l.log_date, l.activities, l.blockers, l.plan_tomorrow, l.evidence_url, l.progress_status, l.submitted_at
		FROM logbooks l
		INNER JOIN mentor_assignments ma ON ma.student_id = l.student_id AND ma.is_active = true
		LEFT JOIN auth_user s ON s.id = l.student_id
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id`,
	Insert: `INSERT INTO logbooks
		(id, student_id, log_date, activities, blockers, plan_tomorrow, evidence_url, progress_status, submitted_at)
		VALUES
		(:id, :student_id, :log_date, :activities, :blockers, :plan_tomorrow, :evidence_url, :progress_status, :submitted_at) RETURNING id`,
	Update: `UPDATE logbooks SET log_date = :log_date, activities = :activities, blockers = :blockers, plan_tomorrow = :plan_tomorrow, evidence_url = :evidence_url, progress_status = :progress_status WHERE id = :id`,
	Count:  `SELECT count(*) FROM logbooks `,
}

type LogbookRepository interface {
	Create(ctx context.Context, data *Logbook) error
	GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []LogbookDTO, err error)
	GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []LogbookDTO, err error)
	ResolveByID(ctx context.Context, id uuid.UUID) (data Logbook, err error)
	Update(ctx context.Context, data Logbook) error
	ResolveAllByStudentID(ctx context.Context, studentID uuid.UUID, req RequestLogbookListFormat) (data pagination.Response, err error)
	ResolveAllByMentorID(ctx context.Context, mentorID uuid.UUID, req RequestLogbookListFormat) (data pagination.Response, err error)
	GetStudentName(ctx context.Context, studentID uuid.UUID) (string, error)
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

func (r *LogbookRepositoryPostgreSQL) GetStudentName(ctx context.Context, studentID uuid.UUID) (string, error) {
	var name string
	err := r.DB.Read.GetContext(ctx, &name, "SELECT name FROM auth_user WHERE id = $1", studentID)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return name, err
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

	if req.ProgressStatus != "" {
		searchBuff.WriteString(fmt.Sprintf(" AND l.progress_status = '%s' ", req.ProgressStatus))
	}

	if req.Date != "" {
		searchBuff.WriteString(fmt.Sprintf(" AND l.log_date = '%s' ", req.Date))
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

	if req.Date != "" {
		// Query for specific date (Monitoring style)
		queryBase := `
			FROM mentor_assignments ma
			INNER JOIN auth_user s ON s.id = ma.student_id
			LEFT JOIN logbooks l ON l.student_id = ma.student_id AND l.log_date::date = ?
			WHERE ma.mentor_id = ? AND ma.is_active = true
		`

		var filterSQL string
		if req.Search != "" {
			filterSQL += fmt.Sprintf(" AND (s.name ILIKE '%%%s%%' OR l.activities ILIKE '%%%s%%' OR l.blockers ILIKE '%%%s%%' OR l.plan_tomorrow ILIKE '%%%s%%') ", req.Search, req.Search, req.Search, req.Search)
		}

		if req.ProgressStatus != "" {
			if req.ProgressStatus == "pending" {
				filterSQL += " AND l.progress_status IS NULL "
			} else {
				filterSQL += fmt.Sprintf(" AND l.progress_status = '%s' ", req.ProgressStatus)
			}
		}

		countQuery := r.DB.Read.Rebind("SELECT count(*) " + queryBase + filterSQL)
		err = r.DB.Read.QueryRowContext(ctx, countQuery, req.Date, mentorID).Scan(&totalData)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}

		if totalData < 1 {
			data.Items = make([]interface{}, 0)
			data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
			return
		}

		offset := (req.PageNumber - 1) * req.PageSize

		selectSQL := `
			SELECT 
				COALESCE(l.id, '00000000-0000-0000-0000-000000000000') AS id,
				ma.student_id,
				s.name AS student_name,
				ma.mentor_id,
				COALESCE(mentor.name, '') AS mentor_name,
				COALESCE(l.log_date, TO_DATE(?, 'YYYY-MM-DD')) AS log_date,
				COALESCE(l.activities, '') AS activities,
				COALESCE(l.blockers, '') AS blockers,
				COALESCE(l.plan_tomorrow, '') AS plan_tomorrow,
				l.evidence_url,
				COALESCE(l.progress_status, 'pending') AS progress_status,
				l.submitted_at
			FROM mentor_assignments ma
			INNER JOIN auth_user s ON s.id = ma.student_id
			LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id
			LEFT JOIN logbooks l ON l.student_id = ma.student_id AND l.log_date::date = ?
			WHERE ma.mentor_id = ? AND ma.is_active = true
		`

		selectSQL += filterSQL
		selectSQL += " ORDER BY s.name ASC LIMIT ? OFFSET ? "

		finalQuery := r.DB.Read.Rebind(selectSQL)
		rows, err := r.DB.Read.QueryxContext(ctx, finalQuery, req.Date, req.Date, mentorID, req.PageSize, offset)
		if err != nil {
			logger.ErrorWithStack(err)
			return data, err
		}
		defer rows.Close()

		for rows.Next() {
			var logbook LogbookDTO
			err = rows.StructScan(&logbook)
			if err != nil {
				logger.ErrorWithStack(err)
				return data, err
			}
			data.Items = append(data.Items, logbook)
		}

		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
		return data, nil
	}

	// Existing behavior (when no date is provided)
	filterBuff.WriteString(" WHERE ma.mentor_id = ? ")

	if req.Search != "" {
		filterBuff.WriteString(fmt.Sprintf(" AND (l.activities ILIKE '%%%s%%' OR l.blockers ILIKE '%%%s%%' OR l.plan_tomorrow ILIKE '%%%s%%' OR s.name ILIKE '%%%s%%') ", req.Search, req.Search, req.Search, req.Search))
	}

	if req.ProgressStatus != "" {
		filterBuff.WriteString(fmt.Sprintf(" AND l.progress_status = '%s' ", req.ProgressStatus))
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
