package internship

import (
	"context"
	"fmt"
	"lms-be/infras"
	"lms-be/shared/logger"
)

type HRDMonitoringRepository interface {
	GetMonitoringData(ctx context.Context, req RequestHRDMonitoringFormat) (data []HRDMonitoringDTO, err error)
}

type HRDMonitoringRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideHRDMonitoringRepositoryPostgreSQL(db *infras.PostgresqlConn) *HRDMonitoringRepositoryPostgreSQL {
	return &HRDMonitoringRepositoryPostgreSQL{DB: db}
}

func (r *HRDMonitoringRepositoryPostgreSQL) GetMonitoringData(ctx context.Context, req RequestHRDMonitoringFormat) (data []HRDMonitoringDTO, err error) {
	query := `
		SELECT 
			u.id AS student_id, 
			u.name AS student_name, 
			COALESCE(reg.university, '-') AS university, 
			mentor.name AS mentor_name,
			lb.status AS logbook_status,
			lb.log_date AS log_date
		FROM auth_user u
		LEFT JOIN internship_registrations reg ON reg.user_id = u.id
		LEFT JOIN mentor_assignments ma ON ma.student_id = u.id AND ma.is_active = true
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id
		LEFT JOIN logbooks lb ON lb.student_id = u.id AND lb.log_date::date = CURRENT_DATE
		WHERE u.role_id = 'HA02' AND u.active = true
	`

	if req.Search != "" {
		query += fmt.Sprintf(" AND (u.name ILIKE '%%%s%%' OR reg.university ILIKE '%%%s%%')", req.Search, req.Search)
	}

	query += " ORDER BY u.name ASC"

	rows, err := r.DB.Read.QueryxContext(ctx, query)
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
