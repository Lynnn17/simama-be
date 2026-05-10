package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type HRDMonitoringDTO struct {
	StudentID        uuid.UUID  `db:"student_id" json:"studentId"`
	StudentName      string     `db:"student_name" json:"studentName"`
	University       string     `db:"university" json:"university"`
	MentorName       *string    `db:"mentor_name" json:"mentorName"`
	AttendanceStatus string     `json:"attendanceStatus"`
	LogbookStatus    *string    `db:"logbook_status" json:"logbookStatus"`
	LogbookDate      *time.Time `db:"log_date" json:"logDate"`
}

func main() {
	dsn := "host=145.79.13.180 port=5432 user=dev password=rahasia123 dbname=magang sslmode=disable"
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	var data []HRDMonitoringDTO
	err = db.Select(&data, query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Results found: %d\n", len(data))
	for _, d := range data {
		fmt.Printf("Student: %s, ID: %s, Active: true\n", d.StudentName, d.StudentID)
	}
}
