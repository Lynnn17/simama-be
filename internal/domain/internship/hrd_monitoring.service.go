package internship

import (
	"context"
	"fmt"
	"time"
)

type HRDMonitoringService interface {
	GetMonitoringData(ctx context.Context, req RequestHRDMonitoringFormat) (data []HRDMonitoringDTO, err error)
	GetStudentQuickView(ctx context.Context, studentID string) (data StudentQuickViewDTO, err error)
}

type HRDMonitoringServiceImpl struct {
	Repo HRDMonitoringRepository
}

func ProvideHRDMonitoringServiceImpl(repo HRDMonitoringRepository) *HRDMonitoringServiceImpl {
	return &HRDMonitoringServiceImpl{Repo: repo}
}

func (s *HRDMonitoringServiceImpl) GetMonitoringData(ctx context.Context, req RequestHRDMonitoringFormat) (data []HRDMonitoringDTO, err error) {
	data, err = s.Repo.GetMonitoringData(ctx, req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("DEBUG: Monitoring Data Length: %d\n", len(data))

	now := time.Now()
	thresholdHour := 17
	todayStr := now.Format("2006-01-02")

	// Parse monitoring date
	monitoringDate := now
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err == nil {
			monitoringDate = parsedDate
		}
	}
	isToday := req.Date == "" || req.Date == todayStr
	weekday := monitoringDate.Weekday()

	// Use WIB timezone for hour comparison
	loc, _ := time.LoadLocation("Asia/Jakarta")
	if loc == nil {
		loc = time.FixedZone("WIB", 7*3600)
	}

	for i := range data {
		// Logika Hari Kerja (Senin-Jumat)
		if weekday != time.Saturday && weekday != time.Sunday {
			if data[i].LogbookStatus != nil {
				// Sudah mengisi logbook (Hadir)
				data[i].AttendanceStatus = "Hadir"

				// Cek jam submit (LogbookDate berisi submitted_at)
				if data[i].LogbookDate != nil {
					wibTime := data[i].LogbookDate.In(loc)
					if wibTime.Hour() < thresholdHour {
						status := "Submitted"
						data[i].LogbookStatus = &status
					} else {
						status := "Late"
						data[i].LogbookStatus = &status
					}
				}
			} else {

				// Belum mengisi logbook
				if isToday && now.In(loc).Hour() < thresholdHour {
					// Masih Jam Kerja (Hari Ini)
					status := "Pending"
					data[i].LogbookStatus = &status
					data[i].AttendanceStatus = "Belum Tercatat"
				} else {
					// Sudah Lewat Jam Kerja atau Hari Kemarin
					status := "Late"
					data[i].LogbookStatus = &status
					data[i].AttendanceStatus = "Tidak Hadir"
				}
			}
		} else {
			// Akhir Pekan
			data[i].AttendanceStatus = "Libur"
			status := "Holiday"
			data[i].LogbookStatus = &status
		}
	}

	return data, nil
}
func (s *HRDMonitoringServiceImpl) GetStudentQuickView(ctx context.Context, studentID string) (data StudentQuickViewDTO, err error) {
	return s.Repo.GetStudentQuickView(ctx, studentID)
}
