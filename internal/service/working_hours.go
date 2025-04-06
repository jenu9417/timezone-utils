package service

import (
	"errors"
	"time"
	"timezone-utils/internal/client"
	"timezone-utils/pkg/models"
)

const workStartHour = 9
const workEndHour = 17

func CheckOverlap(req models.WorkingHoursRequest) (models.WorkingHoursResponse, error) {
	// Parse input time
	t, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return models.WorkingHoursResponse{}, errors.New("invalid datetime format")
	}

	// Convert to both timezones
	loc1, err := time.LoadLocation(req.Timezone1)
	if err != nil {
		return models.WorkingHoursResponse{}, errors.New("invalid timezone 1")
	}
	loc2, err := time.LoadLocation(req.Timezone2)
	if err != nil {
		return models.WorkingHoursResponse{}, errors.New("invalid timezone 2")
	}

	localTime1 := t.In(loc1)
	localTime2 := t.In(loc2)

	isWorkHour1 := isWorkingHour(localTime1)
	isWorkHour2 := isWorkingHour(localTime2)

	isHoliday1, err := client.IsHoliday(req.Country1, localTime1)
	if err != nil {
		isHoliday1 = false // fail-safe
	}

	isHoliday2, err := client.IsHoliday(req.Country2, localTime2)
	if err != nil {
		isHoliday2 = false
	}

	return models.WorkingHoursResponse{
		IsOverlap:  isWorkHour1 && isWorkHour2 && !isHoliday1 && !isHoliday2,
		LocalTime1: localTime1.Format(time.RFC3339),
		LocalTime2: localTime2.Format(time.RFC3339),
		Holiday1:   isHoliday1,
		Holiday2:   isHoliday2,
	}, nil
}

func isWorkingHour(t time.Time) bool {
	hour := t.Hour()
	return hour >= workStartHour && hour < workEndHour
}
