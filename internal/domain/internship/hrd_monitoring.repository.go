package internship

import (
	"context"
	"fmt"
	"lms-be/infras"
	"lms-be/shared/logger"
	"time"
)

type HRDMonitoringRepository interface {
	GetMonitoringData(ctx context.Context, req RequestHRDMonitoringFormat) (data []HRDMonitoringDTO, err error)
	GetStudentQuickView(ctx context.Context, studentID string) (data StudentQuickViewDTO, err error)
}

type HRDMonitoringRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideHRDMonitoringRepositoryPostgreSQL(db *infras.PostgresqlConn) *HRDMonitoringRepositoryPostgreSQL {
	return &HRDMonitoringRepositoryPostgreSQL{DB: db}
}

func (r *HRDMonitoringRepositoryPostgreSQL) GetMonitoringData(ctx context.Context, req RequestHRDMonitoringFormat) (data []HRDMonitoringDTO, err error) {
	monitoringDate := req.Date
	if monitoringDate == "" {
		monitoringDate = time.Now().Format("2006-01-02")
	}

	query := `
		SELECT 
			u.id AS student_id, 
			u.name AS student_name, 
			COALESCE(reg.university, '-') AS university, 
			mentor.name AS mentor_name,
			lb.progress_status AS logbook_status,
			lb.submitted_at AS log_date,
			lb_y.progress_status AS yesterday_logbook_status,
			lb_y.submitted_at AS yesterday_log_date
		FROM auth_user u
		LEFT JOIN internship_registrations reg ON reg.user_id = u.id
		LEFT JOIN mentor_assignments ma ON ma.student_id = u.id AND ma.is_active = true AND ma.assigned_at::date <= $1::date
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id
		LEFT JOIN logbooks lb ON lb.student_id = u.id AND lb.log_date::date = $1
		LEFT JOIN logbooks lb_y ON lb_y.student_id = u.id AND lb_y.log_date::date = (
			CASE 
				WHEN EXTRACT(DOW FROM $1::date) = 1 THEN $1::date - INTERVAL '3 days' 
				WHEN EXTRACT(DOW FROM $1::date) = 0 THEN $1::date - INTERVAL '2 days'
				ELSE $1::date - INTERVAL '1 day' 
			END
		)::date
		WHERE u.role_id = 'HA02' AND u.active = true AND ma.id IS NOT NULL AND ma.is_active = true AND ma.assigned_at::date <= $1::date
	`

	if req.Search != "" {
		query += fmt.Sprintf(" AND (u.name ILIKE '%%%s%%' OR reg.university ILIKE '%%%s%%')", req.Search, req.Search)
	}

	query += " ORDER BY u.name ASC"

	rows, err := r.DB.Read.QueryxContext(ctx, query, monitoringDate)
	if err != nil {
		logger.ErrorWithStack(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d HRDMonitoringDTO
		err = rows.StructScan(&d)
		if err != nil {
			logger.ErrorWithStack(err)
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}
func (r *HRDMonitoringRepositoryPostgreSQL) GetStudentQuickView(ctx context.Context, studentID string) (data StudentQuickViewDTO, err error) {
	query := `
		SELECT 
			u.id AS student_id, 
			u.name AS student_name, 
			COALESCE(reg.university, '-') AS university, 
			COALESCE(reg.major, '-') AS major,
			COALESCE(reg.phone, '-') AS phone,
			COALESCE(reg.email, '-') AS email,
			COALESCE(mentor.name, 'Belum ditugaskan') AS mentor_name,
			(SELECT COUNT(*) FROM logbooks WHERE student_id = $1) AS total_attendance
		FROM auth_user u
		LEFT JOIN internship_registrations reg ON reg.user_id = u.id
		LEFT JOIN mentor_assignments ma ON ma.student_id = u.id AND ma.is_active = true
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id
		WHERE u.id = $1
	`

	err = r.DB.Read.GetContext(ctx, &data, query, studentID)
	if err != nil {
		logger.ErrorWithStack(err)
		return data, err
	}

	// Fetch last logbook
	lastLbQuery := `
		SELECT activities, progress_status AS status, submitted_at
		FROM logbooks
		WHERE student_id = $1
		ORDER BY log_date DESC
		LIMIT 1
	`
	var lastLb LastLogbookDTO
	err = r.DB.Read.GetContext(ctx, &lastLb, lastLbQuery, studentID)
	if err == nil {
		data.LastLogbook = &lastLb
	} else {
		// It's okay if no logbook found
		data.LastLogbook = nil
		err = nil
	}

	return data, nil
}
