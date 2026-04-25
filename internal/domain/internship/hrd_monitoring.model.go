package internship

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type HRDMonitoringDTO struct {
	StudentID        uuid.UUID  `db:"student_id" json:"studentId"`
	StudentName      string     `db:"student_name" json:"studentName"`
	University       string     `db:"university" json:"university"`
	MentorName       *string    `db:"mentor_name" json:"mentorName"`
	AttendanceStatus string     `json:"attendanceStatus"` // Hadir, Tidak Hadir, Belum Tercatat
	LogbookStatus    *string    `db:"logbook_status" json:"logbookStatus"`
	LogbookDate      *time.Time `db:"log_date" json:"logDate"`
}

type RequestHRDMonitoringFormat struct {
	Search string `json:"search"`
}
