package repository

import (
	"database/sql"
	"instructor-led-app/config"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/model"
	"log"
	"math"

	"github.com/sirupsen/logrus"
)

type QuestionRepository interface {
	List(page, size int) ([]entity.Question, model.Paging, error)
	Get(id string) (entity.Question, error)
	Create(payload entity.Question) (entity.Question, error)
	Delete(id string) error
	Update(id string, payload entity.Question) (entity.Question, error)
	GetQuestionByTrainerId(id string, page, size int) ([]entity.Question, model.Paging, error)
	CreateQuestionByTrainer(participantId string, payload dto.QuestionDTO) (dto.QuestionDTO, error)
	UpdateQuestionStatusByTrainer(participantId string, payload dto.QuestionDTO) (dto.QuestionDTO, error)
	GetQuestionByScheduleIdandParticipantId(scheduleId, participantId string) (entity.Question, error)
	CreateQuestionByParticipant(participantId string, payload dto.QuestionDto) (dto.QuestionDto, error)
}

type questionRepository struct {
	db *sql.DB
}

// Create Quetion Boleh gw Nihhhh
func (q *questionRepository) CreateQuestionByParticipant(participantId string, payload dto.QuestionDto) (dto.QuestionDto, error) {
	var question dto.QuestionDto
	status := "Process"

	err := q.db.QueryRow(config.CreateQuestionQuery,
		payload.Question,
		status,
		participantId,
		payload.ScheduleID).Scan(
		&question.ID,
		&question.ParticipantID,
		&question.ScheduleID,
		&question.CreatedAt,
		&question.UpdatedAt)
	if err != nil {
		log.Println("questionRepository.QueryRow:", err.Error())
		return dto.QuestionDto{}, err
	}

	question.Question = payload.Question
	question.ScheduleID = payload.ScheduleID

	return question, nil
}

// GetQuestionByScheduleId implements QuestionRepository.
func (q *questionRepository) GetQuestionByTrainerId(id string, page, size int) ([]entity.Question, model.Paging, error) {
	var questions []entity.Question
	offset := (page - 1) * size
	rows, err := q.db.Query(config.SelectQuestionByTrainerID, id, size, offset)
	if err != nil {
		return nil, model.Paging{}, err
	}

	for rows.Next() {
		var question entity.Question
		err := rows.Scan(&question.ID, &question.Question, &question.Status, &question.ParticipantID, &question.TrainerID, &question.ScheduleID, &question.CreatedAt, &question.UpdatedAt)
		if err != nil {
			return nil, model.Paging{}, err
		}
		questions = append(questions, question)
	}
	totalRows := 0
	if err := q.db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return questions, paging, nil
}

// GetQuestionByScheduleIdandParticipantId implements QuestionRepository.
func (q *questionRepository) GetQuestionByScheduleIdandParticipantId(scheduleId string, participantId string) (entity.Question, error) {
	logger := logrus.New()
	logger.Infoln(scheduleId)
	logger.Infoln(participantId)
	var question entity.Question
	status := "Process"

	err := q.db.QueryRow(config.SelectQuestionByScheduleIDandParticipantID, scheduleId, participantId, status).Scan(
		&question.ID,
		&question.Question,
		&question.Status,
		&question.ParticipantID,
	)
	if err != nil {
		log.Println("questionRepository.QueryRow:", err.Error())
		return entity.Question{}, err
	}
	return question, nil
}

// UpdateQuestionStatusByTrainer implements QuestionRepository.
func (q *questionRepository) UpdateQuestionStatusByTrainer(participantId string, payload dto.QuestionDTO) (dto.QuestionDTO, error) {
	var questionCheck dto.QuestionDTO

	if err := q.db.QueryRow(config.UpdateQuestionStatusByTrainer,
		participantId,
		payload.ScheduleID,
		payload.Answer,
		payload.Status,
		payload.UpdatedAt,
		payload.TrainerID).Scan(
		&questionCheck.ID,
		&questionCheck.Question,
		&questionCheck.Answer,
		&questionCheck.Status,
		&questionCheck.TrainerID,
		&questionCheck.ParticipantID,
		&questionCheck.ScheduleID,
		&questionCheck.CreatedAt,
		&questionCheck.UpdatedAt); err != nil {
		log.Println("QueryRow.err UpdatedByID :", err)
		return dto.QuestionDTO{}, err
	}
	questionCheck.ParticipantName = payload.ParticipantName

	return questionCheck, nil
}

// CreateQuestionByTrainer implements QuestionRepository.
func (q *questionRepository) CreateQuestionByTrainer(participantId string, payload dto.QuestionDTO) (dto.QuestionDTO, error) {
	var question dto.QuestionDTO
	err := q.db.QueryRow(config.InsertQuestionNew,
		payload.Question,
		payload.Answer,
		payload.Status,
		participantId,
		payload.TrainerID,
		payload.ScheduleID).Scan(
		&question.ID,
		&question.ParticipantID,
		&question.TrainerID,
		&question.ScheduleID,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err != nil {
		log.Println("questionRepository.QueryRow:", err.Error())
		return dto.QuestionDTO{}, err
	}

	question.ParticipantName = payload.ParticipantName
	question.Question = payload.Question
	question.Answer = payload.Answer
	question.Status = payload.Status
	question.TrainerID = payload.TrainerID
	question.ScheduleID = payload.ScheduleID

	return question, nil
}

func (q *questionRepository) Create(payload entity.Question) (entity.Question, error) {
	var question entity.Question
	err := q.db.QueryRow(config.InsertQuestion,
		payload.Question,
		payload.Status,
		payload.ParticipantID,
		payload.TrainerID,
		payload.ScheduleID,
		payload.UpdatedAt).Scan(
		&question.ID,
		&question.CreatedAt,
	)
	if err != nil {
		log.Println("questionRepository.QueryRow:", err.Error())
		return entity.Question{}, err
	}
	question.Question = payload.Question
	question.Status = payload.Status
	question.ParticipantID = payload.ParticipantID
	question.TrainerID = payload.TrainerID
	question.ScheduleID = payload.ScheduleID
	question.UpdatedAt = payload.UpdatedAt
	return question, nil
}

func (q *questionRepository) Get(id string) (entity.Question, error) {
	var question entity.Question
	err := q.db.QueryRow(config.SelectQuestionByID, id).Scan(
		&question.ID,
		&question.Question,
		&question.Status,
		&question.ParticipantID,
		&question.TrainerID,
		&question.ScheduleID,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err != nil {
		log.Println("questionRepository.QueryRow:", err.Error())
		return entity.Question{}, err
	}
	return question, nil
}

func (q *questionRepository) Delete(id string) error {
	// var question entity.Question
	result, err := q.db.Exec(config.DeleteQuestion, id)
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

func (q *questionRepository) Update(id string, payload entity.Question) (entity.Question, error) {
	_, err := q.db.Exec(config.UpdateQuestion, id, payload.Status)
	if err != nil {
		log.Println("questionRepository.Exec:", err.Error())
		return entity.Question{}, err
	}
	return payload, nil
}

func (q *questionRepository) List(page, size int) ([]entity.Question, model.Paging, error) {
	var questions []entity.Question
	offset := (page - 1) * size
	rows, err := q.db.Query(config.SelectQuestionList, size, offset)
	if err != nil {
		log.Println("taskRepository.Query:", err.Error())
		return nil, model.Paging{}, err
	}
	for rows.Next() {
		var question entity.Question
		err := rows.Scan(
			&question.ID,
			&question.Question,
			&question.Status,
			&question.ParticipantID,
			&question.TrainerID,
			&question.ScheduleID,
			&question.CreatedAt,
			&question.UpdatedAt)
		if err != nil {
			log.Println("taskRepository.Rows.Next():", err.Error())
			return nil, model.Paging{}, err
		}

		questions = append(questions, question)
	}
	totalRows := 0

	if err := q.db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}
	return questions, paging, nil
}

func NewQuestionRepository(db *sql.DB) QuestionRepository {
	return &questionRepository{db: db}
}
