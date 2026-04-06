package internship

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"lms-be/infras"
	"lms-be/shared/logger"
)

var MentorAssignmentQuery = struct {
	SelectByMentorID  string
	SelectByStudentID string
	Insert            string
	ExistByStudentID  string
}{
	SelectByMentorID: `SELECT ma.id, ma.mentor_id, mentor.name AS mentor_name, ma.student_id, student.name AS student_name
		FROM mentor_assignments ma
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id
		LEFT JOIN auth_user student ON student.id = ma.student_id`,
	SelectByStudentID: `SELECT ma.id, ma.mentor_id, mentor.name AS mentor_name, ma.student_id, student.name AS student_name
		FROM mentor_assignments ma
		LEFT JOIN auth_user mentor ON mentor.id = ma.mentor_id
		LEFT JOIN auth_user student ON student.id = ma.student_id`,
	Insert: `INSERT INTO mentor_assignments (id, mentor_id, student_id)
		VALUES (:id, :mentor_id, :student_id) RETURNING id`,
	ExistByStudentID: `SELECT id FROM mentor_assignments`,
}

type MentorAssignmentRepository interface {
	Create(ctx context.Context, data *MentorAssignment) error
	GetStudentsByMentorID(ctx context.Context, mentorID uuid.UUID) (data []MentorAssignmentDTO, err error)
	ResolveMentorByStudentID(ctx context.Context, studentID uuid.UUID) (data MentorAssignmentDTO, err error)
	ExistByStudentID(ctx context.Context, studentID uuid.UUID) (bool, error)
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
	query := r.DB.Read.Rebind(MentorAssignmentQuery.SelectByMentorID + " WHERE ma.mentor_id = ? ORDER BY student_name ASC")
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
	query := r.DB.Read.Rebind(MentorAssignmentQuery.SelectByStudentID + " WHERE ma.student_id = ?")
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
	query := r.DB.Read.Rebind(MentorAssignmentQuery.ExistByStudentID + " WHERE student_id = ?")
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

func (r *MentorAssignmentRepositoryPostgreSQL) String() string {
	return fmt.Sprintf("%T", r)
}
