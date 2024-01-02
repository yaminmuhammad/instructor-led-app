package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"instructor-led-app/config"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/model"
	"log"
	"math"
	"time"
)

type AbsenceRepository interface {
	List(page, size int) ([]entity.Absence, model.Paging, error)
	ListByDate(startDate, endDate time.Time, page, size int) ([]entity.Absence, model.Paging, error)
	GetAbsencesByParticipantID(id string) (entity.Absence, error)
	GetAbsencesByScheduleID(id string) (dto.AbsenceDTO, error)
	Create(payload []dto.ParticipantScheduleDTO) ([]dto.ParticipantScheduleDTO, error)
	Delete(id string) error
	UpdateByScheduleIDandParticipantID(scheduleId, participantId string, payload dto.AbsenceCheckDTO) (dto.AbsenceCheckDTO, error)
}

type absenceRepository struct {
	db *sql.DB
}

// UpdateByScheduleIDandParticipantID implements AbsenceRepository.
func (a *absenceRepository) UpdateByScheduleIDandParticipantID(scheduleId string, participantId string, payload dto.AbsenceCheckDTO) (dto.AbsenceCheckDTO, error) {
	var absenceCheckDTO dto.AbsenceCheckDTO

	if err := a.db.QueryRow(config.UpdateAbsencesByParticipantName, participantId, scheduleId, payload.Information, payload.Absence_status, payload.Absence_time, payload.Updated_at).Scan(&absenceCheckDTO.ID, &absenceCheckDTO.Information, &absenceCheckDTO.Absence_status, &absenceCheckDTO.Absence_time, &absenceCheckDTO.Updated_at); err != nil {
		log.Println("QueryRow.err UpdatedByID :", err)
		return dto.AbsenceCheckDTO{}, err
	}
	absenceCheckDTO.Name = payload.Name

	return absenceCheckDTO, nil
}

// GetAbsencesByScheduleID implements AbsenceRepository.
func (a *absenceRepository) GetAbsencesByScheduleID(id string) (dto.AbsenceDTO, error) {
	var absence dto.AbsenceDTO
	err := a.db.QueryRow(config.GetAbsencesByScheduleId, id).Scan(
		&absence.ID,
		&absence.Date,
		&absence.Participant_id,
		&absence.Trainer_id,
		&absence.Schedule_id,
	)
	if err != nil {
		return dto.AbsenceDTO{}, err
	}
	return absence, nil
}

// GetAbsencesByParticipantID implements AbsenceRepository.
func (a *absenceRepository) GetAbsencesByParticipantID(id string) (entity.Absence, error) {
	var absence entity.Absence
	err := a.db.QueryRow(config.GetAbsencesById, id).Scan(
		&absence.ID,
		&absence.Date,
		&absence.Information,
		&absence.Absence_status,
		&absence.Absence_time,
		&absence.Created_at,
		&absence.Updated_at,
	)
	if err != nil {
		return entity.Absence{}, err
	}
	return absence, nil
}

// ListByDate implements AbsenceRepository.
func (a *absenceRepository) ListByDate(startDate time.Time, endDate time.Time, page int, size int) ([]entity.Absence, model.Paging, error) {
	var absences []entity.Absence
	offset := (page - 1) * size
	rows, err := a.db.Query(config.ListAbsencebyDate, startDate, endDate, size, offset)
	if err != nil {
		return nil, model.Paging{}, err
	}

	for rows.Next() {
		var absence entity.Absence
		err := rows.Scan(&absence.ID, &absence.Date, &absence.Information, &absence.Absence_status, &absence.Absence_time, &absence.Participant_id, &absence.Created_at, &absence.Updated_at)
		if err != nil {
			return nil, model.Paging{}, err
		}
		absences = append(absences, absence)
	}

	totalRows := 0
	if err := a.db.QueryRow("SELECT COUNT(*) FROM absences").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return absences, paging, nil
}

// Create implements AbsenceRepository.
func (a *absenceRepository) Create(payload []dto.ParticipantScheduleDTO) ([]dto.ParticipantScheduleDTO, error) {
	for _, participantID := range payload {
		if len(participantID.ScheduleID) != len(participantID.Date) {
			fmt.Println("Error: ScheduleID and date slices must have the same length")
			return nil, errors.New("invalid data structure")
		}
		for i := 0; i < len(participantID.ScheduleID); i++ {
			schedule := participantID.ScheduleID[i]
			date := participantID.Date[i]
			_, err := a.db.Exec(config.InsertAbsence,
				date,
				participantID.ID,
				participantID.TrainerID,
				schedule,
			)
			if err != nil {
				fmt.Println("Error inserting data:", err.Error())
				return nil, err

			}
		}
	}
	return payload, nil
}

// Delete implements AbsenceRepository.
func (a *absenceRepository) Delete(id string) error {
	_, err := a.db.Exec(config.DeleteByParticipantId, id)
	if err != nil {
		fmt.Println("error deleting absence:", err)
		return err
	}
	return nil
}

// List implements AbsenceRepository.
func (a *absenceRepository) List(page int, size int) ([]entity.Absence, model.Paging, error) {
	var absences []entity.Absence
	offset := (page - 1) * size
	rows, err := a.db.Query(config.ListAbsence, size, offset)
	if err != nil {
		return nil, model.Paging{}, err
	}

	for rows.Next() {
		var absence entity.Absence
		err := rows.Scan(&absence.ID, &absence.Date, &absence.Information, &absence.Absence_status, &absence.Absence_time, &absence.Participant_id, &absence.Created_at, &absence.Updated_at)
		if err != nil {
			return nil, model.Paging{}, err
		}
		absences = append(absences, absence)
	}

	totalRows := 0
	if err := a.db.QueryRow("SELECT COUNT(*) FROM absences").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return absences, paging, nil
}

func NewAbsenceRepository(db *sql.DB) AbsenceRepository {
	return &absenceRepository{db: db}
}
