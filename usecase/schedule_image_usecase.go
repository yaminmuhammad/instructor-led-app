package usecase

import (
	"fmt"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/repository"
	"sort"
	"time"
)

type ScheduleImageUseCase interface {
	UploadImageActivity(userID, filename, startTimeString, endTimeString string) (dto.ScheduleImagesDTO, error)
}

type scheduleImageUseCase struct {
	scheduleImageRepository repository.ScheduleImageRepository
	scheduleUseCase         ScheduleUseCase
}

func (u *scheduleImageUseCase) UploadImageActivity(userID, filename, startTimeString, endTimeString string) (dto.ScheduleImagesDTO, error) {

	schedules, _, err := u.scheduleUseCase.GetScheduleByTrainerID(userID, 1, 3)
	if err != nil {
		return dto.ScheduleImagesDTO{}, fmt.Errorf("failed to get trainer's schedule: %v", err)
	}

	sort.Slice(schedules, func(i, j int) bool {
		return schedules[i].Date.Before(schedules[j].Date)
	})

	today := time.Now()
	var filteredSchedules []entity.Schedule

	for _, schedule := range schedules {
		if schedule.Date.Before(today) {
			filteredSchedules = append(filteredSchedules, schedule)
		}
	}

	schedule := filteredSchedules[0]

	if startTimeString == "" && endTimeString == "" {
		startTimeString = "19:30"
		endTimeString = "20:30"
	}

	now := time.Now()
	startTime, _ := time.Parse("15:04", startTimeString)
	endTime, _ := time.Parse("15:04", endTimeString)

	hour := now.Hour()
	minute := now.Minute()
	startTimeHour := startTime.Hour()
	startTimeMinute := startTime.Minute()
	endTimeHour := endTime.Hour()
	endTimeMinute := endTime.Minute()

	if (hour < startTimeHour && minute < startTimeMinute) || (hour > endTimeHour && minute > endTimeMinute) {
		return dto.ScheduleImagesDTO{}, fmt.Errorf("cannot upload image outside the schedule time")
	}

	imageDTO, err := u.scheduleImageRepository.Insert(dto.ScheduleImagesDTO{
		ScheduleID: schedule.ID,
		FileName:   filename,
	})
	if err != nil {
		return dto.ScheduleImagesDTO{}, fmt.Errorf("failed to save image: %v", err)
	}

	return imageDTO, nil
}

func NewScheduleImageUseCase(scheduleImageRepository repository.ScheduleImageRepository, scheduleUseCase ScheduleUseCase) ScheduleImageUseCase {
	return &scheduleImageUseCase{scheduleImageRepository, scheduleUseCase}
}
