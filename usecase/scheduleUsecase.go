package usecase

import (
	"fmt"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/repository"
	"instructor-led-app/shared/model"
	"time"
)

type ScheduleUseCase interface {
	InsertNewSchedule(payload dto.ScheduleDto) (dto.ScheduleDto, error)
	FindAllSchedule(startDate, endDate time.Time, page, size int) ([]entity.Schedule, model.Paging, error)
	GetScheduleByParticipantID(id string, page, size int) ([]entity.Schedule, model.Paging, error)
	GetScheduleWithParticipantId(id string) ([]string, error)
	GetScheduleByTrainerID(userID string, page, size int) ([]entity.Schedule, model.Paging, error)
	UpdateScheduleByAdmin(trainerId string, code int) ([]entity.Schedule, error)
	DeleteScheduleByDate(date string) error
}

type scheduleUseCase struct {
	repo           repository.ScheduleRepository
	trainerUseCase TrainerUsecase
}

// GetScheduleWithParticipantId implements ScheduleUseCase.
func (s *scheduleUseCase) GetScheduleWithParticipantId(id string) ([]string, error) {
	return s.repo.GetScheduleByTrainerIdWithoutPagination(id)
}

// DeleteScheduleByID implements ScheduleUseCase.
func (s *scheduleUseCase) DeleteScheduleByDate(date string) error {
	return s.repo.DeleteByDate(date)
}

// UpdateScheduleByAdmin implements ScheduleUseCase.
func (s *scheduleUseCase) UpdateScheduleByAdmin(trainerId string, code int) ([]entity.Schedule, error) {
	date, err := s.repo.GetDateByDay(code)
	if err != nil {
		return []entity.Schedule{}, fmt.Errorf("failed to get date: %v", err.Error())
	}

	update, err := s.repo.UpdateScheduleByAdmin(trainerId, date)
	if err != nil {
		return []entity.Schedule{}, fmt.Errorf("failed to update :%v", err.Error())
	}

	return update, nil
}

// FindAllSchedule implements ScheduleUseCase.
func (s *scheduleUseCase) FindAllSchedule(startDate time.Time, endDate time.Time, page int, size int) ([]entity.Schedule, model.Paging, error) {
	if !startDate.IsZero() && !endDate.IsZero() {
		return s.repo.ListByDate(startDate, endDate, page, size)
	} else {
		return s.repo.List(page, size)
	}
}

// GetScheduleByParticipantID implements ScheduleUseCase.
func (s *scheduleUseCase) GetScheduleByParticipantID(id string, page int, size int) ([]entity.Schedule, model.Paging, error) {
	return s.repo.GetScheduleByParticipantID(id, page, size)
}

// GetScheduleByTrainerID implements ScheduleUseCase.
func (s *scheduleUseCase) GetScheduleByTrainerID(userID string, page int, size int) ([]entity.Schedule, model.Paging, error) {
	trainer, err := s.trainerUseCase.FindTrainerByUserId(userID)
	if err != nil {
		return []entity.Schedule{}, model.Paging{}, fmt.Errorf("failed to get trainer: %v", err)
	}

	schedule, paging, err := s.repo.GetScheduleByTrainerId(trainer.ID, page, size)
	if err != nil {
		return []entity.Schedule{}, model.Paging{}, fmt.Errorf("failed to get schedule's : %s", err.Error())

	}

	if len(schedule) == 0 {
		return []entity.Schedule{}, model.Paging{}, fmt.Errorf("no schedule found for the trainer")
	}
	return schedule, paging, nil
}

// InsertNewSchedule implements ScheduleUseCase.
func (s *scheduleUseCase) InsertNewSchedule(payload dto.ScheduleDto) (dto.ScheduleDto, error) {
	if payload.Activity == "" || payload.ParticipantID == "" || payload.TrainerID == "" {
		return dto.ScheduleDto{}, fmt.Errorf("oops, Required field is empty")
	}
	payload.Date = time.Now().Format("2006-01-02")
	schedule, err := s.repo.Create(payload)
	if err != nil {
		return dto.ScheduleDto{}, fmt.Errorf("oppps, failed to save data absence :%v", err.Error())
	}
	return schedule, nil
}

func NewScheduleUseCase(repo repository.ScheduleRepository, trainerUsecae TrainerUsecase) ScheduleUseCase {
	return &scheduleUseCase{repo: repo, trainerUseCase: trainerUsecae}
}
