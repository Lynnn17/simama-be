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
	AttendanceStatus         string     `json:"attendanceStatus"` // Hadir, Tidak Hadir, Belum Tercatat
	LogbookStatus            *string    `db:"logbook_status" json:"logbookStatus"`
	LogbookDate              *time.Time `db:"log_date" json:"logDate"`
	YesterdayLogbookStatus   *string    `db:"yesterday_logbook_status" json:"yesterdayLogbookStatus"`
	YesterdayLogDate         *time.Time `db:"yesterday_log_date" json:"yesterdayLogDate"`
}

type RequestHRDMonitoringFormat struct {
	Search string `json:"search"`
	Date   string `json:"date"`
}

type StudentQuickViewDTO struct {
	StudentID       uuid.UUID       `db:"student_id" json:"studentId"`
	StudentName     string          `db:"student_name" json:"studentName"`
	University      string          `db:"university" json:"university"`
	Major           string          `db:"major" json:"major"`
	Phone           string          `db:"phone" json:"phone"`
	Email           string          `db:"email" json:"email"`
	MentorName      string          `db:"mentor_name" json:"mentorName"`
	TotalAttendance int             `db:"total_attendance" json:"totalAttendance"`
	LastLogbook     *LastLogbookDTO `json:"lastLogbook"`
}

type LastLogbookDTO struct {
	Activities  string    `db:"activities" json:"activities"`
	Status      string    `db:"status" json:"status"`
	SubmittedAt time.Time `db:"submitted_at" json:"submittedAt"`
}
