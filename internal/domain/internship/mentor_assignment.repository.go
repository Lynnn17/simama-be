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

var MentorAssignmentQuery = struct {
	SelectByMentorID  string
	SelectByStudentID string
	Insert            string
	ExistByStudentID  string
	Count             string
	Update            string
}{
	SelectByMentorID: `SELECT ma.id, ma.mentor_id, mentor.name AS mentor_name, ma.student_id, student.name AS student_name, ma.assigned_by, ma.assigned_at, ma.is_active
		FROM mentor_assignments ma
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id
		LEFT JOIN auth_user student ON student.id = ma.student_id`,
	SelectByStudentID: `SELECT ma.id, ma.mentor_id, mentor.name AS mentor_name, ma.student_id, student.name AS student_name, ma.assigned_by, ma.assigned_at, ma.is_active
		FROM mentor_assignments ma
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id
		LEFT JOIN auth_user student ON student.id = ma.student_id`,
	Insert: `INSERT INTO mentor_assignments (id, mentor_id, student_id, assigned_by, assigned_at, is_active)
		VALUES (:id, :mentor_id, :student_id, :assigned_by, :assigned_at, :is_active) RETURNING id`,
	ExistByStudentID: `SELECT id FROM mentor_assignments`,
	Count:            `SELECT count(id) FROM mentor_assignments`,
	Update:           `UPDATE mentor_assignments SET mentor_id = :mentor_id, student_id = :student_id, is_active = :is_active WHERE id = :id`,
}

type MentorAssignmentRepository interface {
	Create(ctx context.Context, data *MentorAssignment) error
	GetStudentsByMentorID(ctx context.Context, mentorID uuid.UUID) (data []MentorAssignmentDTO, err error)
	ResolveMentorByStudentID(ctx context.Context, studentID uuid.UUID) (data MentorAssignmentDTO, err error)
	ExistByStudentID(ctx context.Context, studentID uuid.UUID) (bool, error)
	ResolveAll(ctx context.Context, req RequestMentorAssignmentListFormat) (data pagination.Response, err error)
	Update(ctx context.Context, data *MentorAssignment) error
	DeactivateByStudentID(ctx context.Context, studentID uuid.UUID) error
}

type MentorAssignmentRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideMentorAssignmentRepositoryPostgreSQL(db *infras.PostgresqlConn) *MentorAssignmentRepositoryPostgreSQL {
	s := new(MentorAssignmentRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *MentorAssignmentRepositoryPostgreSQL) Create(ctx context.Context, data *MentorAssignment) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, MentorAssignmentQuery.Insert)
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

func (r *MentorAssignmentRepositoryPostgreSQL) GetStudentsByMentorID(ctx context.Context, mentorID uuid.UUID) (data []MentorAssignmentDTO, err error) {
	query := r.DB.Read.Rebind(MentorAssignmentQuery.SelectByMentorID + " WHERE ma.mentor_id = ? AND ma.is_active = true ORDER BY student_name ASC")
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
		var assignment MentorAssignmentDTO
		err = rows.StructScan(&assignment)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data = append(data, assignment)
	}
	return
}

func (r *MentorAssignmentRepositoryPostgreSQL) ResolveMentorByStudentID(ctx context.Context, studentID uuid.UUID) (data MentorAssignmentDTO, err error) {
	query := r.DB.Read.Rebind(MentorAssignmentQuery.SelectByStudentID + " WHERE ma.student_id = ? AND ma.is_active = true")
	err = r.DB.Read.GetContext(ctx, &data, query, studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return MentorAssignmentDTO{}, err
		}
		logger.ErrorWithStack(err)
		return MentorAssignmentDTO{}, err
	}
	return
}

func (r *MentorAssignmentRepositoryPostgreSQL) ExistByStudentID(ctx context.Context, studentID uuid.UUID) (bool, error) {
	var idResult uuid.UUID
	query := r.DB.Read.Rebind(MentorAssignmentQuery.ExistByStudentID + " WHERE student_id = ? AND is_active = true")
	err := r.DB.Read.GetContext(ctx, &idResult, query, studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.ErrorWithStack(err)
		return false, err
	}
	return true, nil
}

func (r *MentorAssignmentRepositoryPostgreSQL) ResolveAll(ctx context.Context, req RequestMentorAssignmentListFormat) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchBuff bytes.Buffer
	searchBuff.WriteString(" WHERE 1=1 ")

	if req.MentorID != nil && *req.MentorID != uuid.Nil {
		searchBuff.WriteString(" AND ma.mentor_id = ? ")
		searchParams = append(searchParams, req.MentorID)
	}

	if req.StudentID != nil && *req.StudentID != uuid.Nil {
		searchBuff.WriteString(" AND ma.student_id = ? ")
		searchParams = append(searchParams, req.StudentID)
	}

	if req.IsActive != nil {
		searchBuff.WriteString(" AND ma.is_active = ? ")
		searchParams = append(searchParams, req.IsActive)
	}

	if req.Search != "" {
		searchBuff.WriteString(" AND (mentor.name ILIKE ? OR student.name ILIKE ?) ")
		searchParams = append(searchParams, "%"+req.Search+"%", "%"+req.Search+"%")
	}

	var totalData int
	countQuery := r.DB.Read.Rebind("SELECT count(*) FROM mentor_assignments ma " +
		" LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id " +
		" LEFT JOIN auth_user student ON student.id = ma.student_id " +
		searchBuff.String())
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

	searchBuff.WriteString(" ORDER BY ma.assigned_at DESC ")
	searchBuff.WriteString(" LIMIT ? OFFSET ? ")

	offset := (req.PageNumber - 1) * req.PageSize
	query := r.DB.Read.Rebind(MentorAssignmentQuery.SelectByMentorID + searchBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, append(searchParams, req.PageSize, offset)...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var assignment MentorAssignmentDTO
		err = rows.StructScan(&assignment)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		data.Items = append(data.Items, assignment)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *MentorAssignmentRepositoryPostgreSQL) Update(ctx context.Context, data *MentorAssignment) error {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, MentorAssignmentQuery.Update)
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

func (r *MentorAssignmentRepositoryPostgreSQL) DeactivateByStudentID(ctx context.Context, studentID uuid.UUID) error {
	query := r.DB.Write.Rebind("UPDATE mentor_assignments SET is_active = false WHERE student_id = ? AND is_active = true")
	_, err := r.DB.Write.ExecContext(ctx, query, studentID)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return err
}

func (r *MentorAssignmentRepositoryPostgreSQL) String() string {
	return fmt.Sprintf("%T", r)
}
