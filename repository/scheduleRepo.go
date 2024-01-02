package repository

import (
	"database/sql"
	"fmt"
	"instructor-led-app/config"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/model"
	"log"
	"math"
	"time"
)

type ScheduleRepository interface {
	List(page, size int) ([]entity.Schedule, model.Paging, error)
	GetScheduleByTrainerId(id string, page, size int) ([]entity.Schedule, model.Paging, error)
	GetScheduleByTrainerIdWithoutPagination(id string) ([]string, error)
	GetScheduleByParticipantID(id string, page, size int) ([]entity.Schedule, model.Paging, error)
	Create(payload dto.ScheduleDto) (dto.ScheduleDto, error)
	ListByDate(startDate, endDate time.Time, page, size int) ([]entity.Schedule, model.Paging, error)
	GetScheduleIdByDay(code int) ([]string, error)
	GetDateByDay(code int) ([]time.Time, error)
	GetDateByScheduleId(id string) (time.Time, error)
	GetScheduleIdByDate(date time.Time) (string, error)
	GetScheduleWithParticipantId(id string) ([]dto.ScheduleDto, error)
	UpdateScheduleByAdmin(trainerId string, dates []time.Time) ([]entity.Schedule, error)
	DeleteByDate(date string) error
}

type scheduleRepository struct {
	db *sql.DB
}

// DeleteByDate implements ScheduleRepository.
func (s *scheduleRepository) DeleteByDate(date string) error {
	// var question entity.Question
	result, err := s.db.Exec(config.DeleteSchedule, date)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		} else {
			log.Println("questionRepository.Exec:", err.Error())
			return err
		}
	}
	rowAffected, err := result.RowsAffected()
	if rowAffected == 0 {
		return err
	}
	return nil
}

// updateScheduleByAdmin implements ScheduleRepository.
func (s *scheduleRepository) UpdateScheduleByAdmin(trainerId string, dates []time.Time) ([]entity.Schedule, error) {
	var updatedSchedules []entity.Schedule

	for _, date := range dates {
		var updatedSchedule entity.Schedule
		err := s.db.QueryRow(config.UpdateScheduleByAdmin, date, trainerId).Scan(&updatedSchedule.ID, &updatedSchedule.Activity, &updatedSchedule.Date, &updatedSchedule.TrainerID, &updatedSchedule.ParticipantID, &updatedSchedule.UpdatedAt)
		if err != nil {
			log.Println("QueryRow.err UpdatedByID :", err)
			return nil, err
		}
		updatedSchedules = append(updatedSchedules, updatedSchedule)
	}

	return updatedSchedules, nil
}

// GetScheduleIdByDate implements ScheduleRepository.
func (s *scheduleRepository) GetScheduleIdByDate(date time.Time) (string, error) {
	var id string
	err := s.db.QueryRow(config.ScheduleIDbyDate, date).Scan(
		&id,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetDateByScheduleId implements ScheduleRepository.
func (s *scheduleRepository) GetDateByScheduleId(id string) (time.Time, error) {
	var date time.Time
	err := s.db.QueryRow(config.DateByScheduleId, id).Scan(
		&date,
	)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

// GetScheduleByTrainerIdWithoutPagination implements ScheduleRepository.
func (s *scheduleRepository) GetScheduleByTrainerIdWithoutPagination(id string) ([]string, error) {
	var schedules []string
	rows, err := s.db.Query(config.ScheduleIDByTrainerId, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var schedule string
		err := rows.Scan(&schedule)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (s *scheduleRepository) GetScheduleWithParticipantId(id string) ([]dto.ScheduleDto, error) {
	var schedules []dto.ScheduleDto
	rows, err := s.db.Query(config.ScheduleIDByParticipantId, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var schedule dto.ScheduleDto
		if err := rows.Scan(&schedule.ID); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

// GetDateByDay implements ScheduleRepository.
func (s *scheduleRepository) GetDateByDay(code int) ([]time.Time, error) {
	var schedules []time.Time
	rows, err := s.db.Query(config.DateByDay, code)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var schedule time.Time
		err := rows.Scan(&schedule)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// GetScheduleIdByDay implements ScheduleRepository.
func (s *scheduleRepository) GetScheduleIdByDay(code int) ([]string, error) {
	var schedules []string
	rows, err := s.db.Query(config.ScheduleIdByDay, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var schedule string
		err := rows.Scan(&schedule)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

// ListByDate implements ScheduleRepository.
func (s *scheduleRepository) ListByDate(startDate time.Time, endDate time.Time, page int, size int) ([]entity.Schedule, model.Paging, error) {
	var schedules []entity.Schedule
	offset := (page - 1) * size
	rows, err := s.db.Query(config.ListScheduleByDate, startDate, endDate, size, offset)
	if err != nil {
		return nil, model.Paging{}, err
	}
	for rows.Next() {
		var schedule entity.Schedule
		err := rows.Scan(&schedule.ID, &schedule.Activity, &schedule.Date, &schedule.TrainerID, &schedule.ParticipantID, &schedule.CreatedAt, &schedule.UpdatedAt)
		if err != nil {
			return nil, model.Paging{}, err
		}
		schedules = append(schedules, schedule)
	}
	totalRows := 0
	if err := s.db.QueryRow("SELECT COUNT(*) FROM schedules").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return schedules, paging, nil
}

// Create implements ScheduleRepository.
func (s *scheduleRepository) Create(payload dto.ScheduleDto) (dto.ScheduleDto, error) {
	var schedule dto.ScheduleDto
	err := s.db.QueryRow(config.InsertSchedule,
		payload.Activity,
		payload.Date,
		payload.TrainerID,
		payload.ParticipantID).Scan(
		&schedule.ID,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		fmt.Println("task repo error :", err.Error())
		return dto.ScheduleDto{}, err
	}

	schedule.Date = payload.Date
	schedule.Activity = payload.Activity
	schedule.TrainerID = payload.TrainerID
	schedule.ParticipantID = payload.ParticipantID
	return schedule, nil
}

func (s *scheduleRepository) GetScheduleByParticipantID(id string, page, size int) ([]entity.Schedule, model.Paging, error) {
	var schedules []entity.Schedule
	offset := (page - 1) * size
	rows, err := s.db.Query(config.ListScheduleByParticipantId, id, size, offset)
	if err != nil {
		return nil, model.Paging{}, err
	}

	for rows.Next() {
		var schedule entity.Schedule
		err := rows.Scan(&schedule.ID, &schedule.Activity, &schedule.Date, &schedule.TrainerID, &schedule.ParticipantID, &schedule.CreatedAt, &schedule.UpdatedAt)
		if err != nil {
			return nil, model.Paging{}, err
		}
		schedules = append(schedules, schedule)
	}
	totalRows := 0
	if err := s.db.QueryRow("SELECT COUNT(*) FROM schedules").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return schedules, paging, nil

}

// GetScheduleByTrainerId implements ScheduleRepository.
func (s *scheduleRepository) GetScheduleByTrainerId(id string, page, size int) ([]entity.Schedule, model.Paging, error) {
	var schedules []entity.Schedule
	offset := (page - 1) * size
	rows, err := s.db.Query(config.ListScheduleByTrainerId, id, size, offset)
	if err != nil {
		return nil, model.Paging{}, err
	}

	for rows.Next() {
		var schedule entity.Schedule
		err := rows.Scan(&schedule.ID, &schedule.Activity, &schedule.Date, &schedule.TrainerID, &schedule.ParticipantID, &schedule.CreatedAt, &schedule.UpdatedAt)
		if err != nil {
			return nil, model.Paging{}, err
		}
		schedules = append(schedules, schedule)
	}
	totalRows := 0
	if err := s.db.QueryRow("SELECT COUNT(*) FROM schedules").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return schedules, paging, nil
}

// List implements ScheduleRepository.
func (s *scheduleRepository) List(page int, size int) ([]entity.Schedule, model.Paging, error) {
	var schedules []entity.Schedule
	offset := (page - 1) * size
	rows, err := s.db.Query(config.ListSchedule, size, offset)
	if err != nil {
		return nil, model.Paging{}, err
	}
	for rows.Next() {
		var schedule entity.Schedule
		err := rows.Scan(&schedule.ID, &schedule.Activity, &schedule.Date, &schedule.TrainerID, &schedule.ParticipantID, &schedule.CreatedAt, &schedule.UpdatedAt)
		if err != nil {
			return nil, model.Paging{}, err
		}
		schedules = append(schedules, schedule)
	}

	totalRows := 0
	if err := s.db.QueryRow("SELECT COUNT(*) FROM schedules").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return schedules, paging, nil
}

func NewScheduleRepository(db *sql.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}
