package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"instructor-led-app/config"
	"time"

	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/model"
	"log"
	"math"
)

type ParticipantRepository interface {
	Insert(participant dto.ParticipantDTO) (dto.ParticipantDTO, error)
	FindAll(page, size int) ([]dto.ParticipantDTO, model.Paging, error)
	FindByID(id string) (dto.ParticipantDTO, error)
	UpdateByID(participantDto dto.ParticipantDTO, updatedAt string) (dto.ParticipantDTO, error)
	DeleteByID(id string) error
	FindAllIDandRole() ([]dto.ParticipantIdDTO, error)
	GetParticipantByUserId(id string) (dto.ParticipantDTO, error)
	GetScheduleById(participantId string) ([]dto.ParticipantDTO, error)
	UpdateByRole(role, id string) error
}

type participantRepository struct {
	db *sql.DB
}

// UpdateByRole implements ParticipantRepository.
func (r *participantRepository) UpdateByRole(role string, id string) error {
	r.db.Exec(config.UpdateParticipantByRole, id, role)
	return nil
}

// GetParticipantByUserId implements ParticipantRepository.
func (r *participantRepository) GetParticipantByUserId(id string) (dto.ParticipantDTO, error) {
	var participant dto.ParticipantDTO

	if err := r.db.QueryRow(config.GetParticipantBUserID, id).Scan(
		&participant.ID,
		&participant.DateOfBirth,
		&participant.PlaceOfBirth,
		&participant.LastEducation,
		&participant.UserID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.ParticipantDTO{}, fmt.Errorf("participant with id '%s' not found", id)
		}

		return dto.ParticipantDTO{}, err
	}

	return participant, nil
}

func (r *participantRepository) SaveQuestion(questionDTO dto.QuestionDTO) error {
	_, err := r.db.Exec("INSERT INTO questions (question, participant_id, trainer_id, schedule_id, created_at) VALUES ($1, $2, $3, $4, NOW())",
		questionDTO.Question, questionDTO.ParticipantID, questionDTO.TrainerID, questionDTO.ScheduleID)
	// foreign key schedule id dimatikan
	if err != nil {
		return err
	}

	return nil
}

// FindAllIDandRole implements ParticipantRepository.
func (r *participantRepository) FindAllIDandRole() ([]dto.ParticipantIdDTO, error) {

	var participants []dto.ParticipantIdDTO
	rows, err := r.db.Query(config.ListParticipantByIDandRole)
	if err != nil {
		log.Println("List participant error :", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		participant := dto.ParticipantIdDTO{}
		if err := rows.Scan(
			&participant.ID,
			&participant.Role,
		); err != nil {
			log.Println("Scan List participant error :", err)
			return nil, err
		}

		participants = append(participants, participant)
	}

	return participants, nil
}

func (r *participantRepository) CreateNewQuestion(participantId, question string) (dto.ParticipantDTO, error) {
	var participantDTO dto.ParticipantDTO

	rows, err := r.db.Query(config.GetScheduleByParticipantID, participantId)
	if err != nil {
		return participantDTO, err
	}
	defer rows.Close()

	for rows.Next() {
		var scheduleId, trainerId string
		err := rows.Scan(&trainerId, &scheduleId)
		if err != nil {
			log.Println("ParticipantsRepository.ParticipantID.Query:", err.Error())
			return participantDTO, err
		}

		questionDTO := dto.QuestionDTO{
			ParticipantID: participantId,
			ScheduleID:    scheduleId,
			TrainerID:     trainerId,
			Question:      question,
		}

		err = r.SaveQuestion(questionDTO)
		if err != nil {
			log.Println("ParticipantsRepository.question.Query:", err.Error())
			return participantDTO, err
		}

		participantDTO.Schedules = append(participantDTO.Schedules, dto.ScheduleDto{
			ID:        scheduleId,
			TrainerID: trainerId,
			Question: []dto.QuestionDTO{
				{
					ID:            questionDTO.ID,
					ParticipantID: participantId,
					Question:      questionDTO.Question,
					TrainerID:     trainerId,
					ScheduleID:    scheduleId,
				},
			},
		})
	}
	return participantDTO, nil
}

// GetScheduleById implements ParticipantRepository.
func (p *participantRepository) GetScheduleById(id string) ([]dto.ParticipantDTO, error) {
	var participants []dto.ParticipantDTO
	var rows *sql.Rows
	var err error

	auth, err := p.FindByID(id)
	if err != nil {
		log.Println("Error fetching id:", err)
		return nil, err
	}

	if auth.Role == "Basic" {
		// cek hari
		if time.Now().Weekday() == time.Tuesday {
			// query basic
			rows, err = p.db.Query(config.SelectParticipantRoleBasic, id)
		} else {
			return nil, errors.New("tidak ada jadwal")
		}
	} else {
		if time.Now().Weekday() == time.Tuesday {
			rows, err = p.db.Query(config.SelectParticipantRoleAdvance, id)
		} else {
			return nil, errors.New("tidak ada jadwal")
		}
	}

	if err != nil {
		log.Println("ParticipantsRepository.Role.Query:", err.Error())
		return nil, err
	}

	for rows.Next() {
		var participant dto.ParticipantDTO
		err := rows.Scan(&participant.ID, &participant.DateOfBirth, &participant.PlaceOfBirth, &participant.LastEducation, &participant.UserID, &participant.Role)
		if err != nil {
			log.Println("ParticipantsRepository.Rows.Next:", err.Error())
			return nil, err
		}

		participantRows, err := p.db.Query(config.SelectParticipantWithSchedule, participant.ID)
		if err != nil {
			log.Println("ParticipantsRepository.Get.Query:", err.Error())
			return nil, err
		}
		for participantRows.Next() {
			var schedule dto.ScheduleDto
			err := participantRows.Scan(&schedule.ID, &schedule.Activity, &schedule.Day, &schedule.Date, &schedule.TrainerID, &schedule.CreatedAt, &schedule.UpdatedAt)
			if err != nil {
				log.Println("ScheduleRepository.scheduleRows.Next():", err.Error())
				return nil, err
			}
			participant.Schedules = append(participant.Schedules, schedule)
		}
		participants = append(participants, participant)
	}
	return participants, nil
}

func (r *participantRepository) Insert(participant dto.ParticipantDTO) (dto.ParticipantDTO, error) {
	var participantDto dto.ParticipantDTO

	if err := r.db.QueryRow(config.InsertParticipant, participant.DateOfBirth, participant.PlaceOfBirth, participant.LastEducation, participant.UserID, participant.Role).Scan(
		&participantDto.ID,
		&participantDto.DateOfBirth,
		&participantDto.PlaceOfBirth,
		&participantDto.LastEducation,
		&participantDto.UserID,
		&participantDto.Role,
	); err != nil {
		log.Println("Scan error insert participant :", err)
		return dto.ParticipantDTO{}, err
	}

	return participantDto, nil

}

func (r *participantRepository) FindAll(page, size int) ([]dto.ParticipantDTO, model.Paging, error) {
	var participants []dto.ParticipantDTO

	if page == 0 || size == 0 {
		page = 1
		size = 5
	}

	var offset = (page - 1) * size

	rows, err := r.db.Query(config.ListParticipant, size, offset)
	if err != nil {
		log.Println("List participant error :", err)
		return nil, model.Paging{}, err
	}
	defer rows.Close()

	for rows.Next() {
		participant := dto.ParticipantDTO{}
		if err := rows.Scan(
			&participant.ID,
			&participant.DateOfBirth,
			&participant.PlaceOfBirth,
			&participant.LastEducation,
			&participant.UserID,
			&participant.Role,
		); err != nil {
			log.Println("Scan List participant error :", err)
			return nil, model.Paging{}, err
		}

		participants = append(participants, participant)
	}

	totalRows := 0
	if err := r.db.QueryRow("SELECT COUNT(*) FROM participants").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return participants, paging, nil
}

func (r *participantRepository) FindByID(id string) (dto.ParticipantDTO, error) {
	var participant dto.ParticipantDTO

	if err := r.db.QueryRow(config.GetParticipantByID, id).Scan(
		&participant.ID,
		&participant.DateOfBirth,
		&participant.PlaceOfBirth,
		&participant.LastEducation,
		&participant.UserID,
		&participant.Role,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.ParticipantDTO{}, fmt.Errorf("participant with id '%s' not found", id)
		}

		return dto.ParticipantDTO{}, err
	}

	return participant, nil
}

func (r *participantRepository) UpdateByID(participantDto dto.ParticipantDTO, updatedAt string) (dto.ParticipantDTO, error) {
	var participant dto.ParticipantDTO
	if err := r.db.QueryRow(config.UpdateParticipantByID, participantDto.ID, participantDto.DateOfBirth, participantDto.PlaceOfBirth, participantDto.LastEducation, participantDto.UserID, participantDto.Role, updatedAt).Scan(&participant.ID, &participant.DateOfBirth, &participant.PlaceOfBirth, &participant.LastEducation, &participant.UserID, &participant.Role); err != nil {
		log.Println("QueryRow.err UpdatedByID :", err)
		return dto.ParticipantDTO{}, err
	}

	return participant, nil
}

func (r *participantRepository) DeleteByID(id string) error {
	result, err := r.db.Exec(config.DeleteParticipantByID, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("participant with id '%s' not found", id)
	}

	return nil
}

func NewParticipantRepository(db *sql.DB) ParticipantRepository {
	return &participantRepository{db: db}
}
