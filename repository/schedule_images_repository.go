package repository

import (
	"database/sql"
	"instructor-led-app/config"
	"instructor-led-app/entity/dto"
)

type ScheduleImageRepository interface {
	Insert(scheduleImage dto.ScheduleImagesDTO) (dto.ScheduleImagesDTO, error)
}

type scheduleImageRepository struct {
	db *sql.DB
}

func (r *scheduleImageRepository) Insert(scheduleImage dto.ScheduleImagesDTO) (dto.ScheduleImagesDTO, error) {
	var scheduleImageDTO dto.ScheduleImagesDTO

	if err := r.db.QueryRow(config.InsertScheduleImage, scheduleImage.ScheduleID, scheduleImage.FileName).Scan(&scheduleImageDTO.ID, &scheduleImageDTO.ScheduleID, &scheduleImageDTO.FileName); err != nil {
		return dto.ScheduleImagesDTO{}, err
	}

	return scheduleImageDTO, nil
}

func NewScheduleImagesRepository(db *sql.DB) ScheduleImageRepository {
	return &scheduleImageRepository{db}
}
