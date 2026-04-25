package internship

import (
	"context"
	"time"
)

type HRDMonitoringService interface {
	GetMonitoringData(ctx context.Context, req RequestHRDMonitoringFormat) (data []HRDMonitoringDTO, err error)
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

	now := time.Now()
	// Set threshold jam 17:00
	thresholdHour := 17

	for i := range data {
		// Logika Kehadiran
		if data[i].LogbookStatus != nil {
			data[i].AttendanceStatus = "Hadir"
		} else {
			if now.Hour() < thresholdHour {
				data[i].AttendanceStatus = "Belum Tercatat"
			} else {
				data[i].AttendanceStatus = "Tidak Hadir"
			}
		}
	}

	return data, nil
}
