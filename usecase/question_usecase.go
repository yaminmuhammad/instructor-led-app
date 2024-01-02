package usecase

import (
	"fmt"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/repository"
	"instructor-led-app/shared/model"
	"time"

	"github.com/sirupsen/logrus"
)

type QuestionUseCase interface {
	FindById(id string) (entity.Question, error)
	FindAllQuestion(page, size int) ([]entity.Question, model.Paging, error)
	DeleteQuestion(id string) error
	CreateNewQuestion(payload entity.Question) (entity.Question, error)
	UpdateQuestion(id string, payload entity.Question) (entity.Question, error)
	FindQuestionByTrainerId(userID string, page, size int) ([]entity.Question, model.Paging, error)
	CreateQuestionByTrainer(trainerId, participantId string, payload dto.QuestionDTO) (dto.QuestionDTO, error)
	UpdatadStatusQuestionByTrainer(trainerId, participantId string, payload dto.QuestionDTO) (dto.QuestionDTO, error)
	CreateQuestionByParticipant(participantId string, payload dto.QuestionDto) (dto.QuestionDto, error)
}

type questionUseCase struct {
	repo            repository.QuestionRepository
	participantRepo repository.ParticipantRepository
	scheduleRepo    repository.ScheduleRepository
	userRepo        repository.UserRepository
	trainerRepo     repository.TrainerRepository
	participantUC   ParticipantUseCase
	trainerUseCase  TrainerUsecase
}

// CreateQuestionByParticipant implements QuestionUseCase.

// UpdatadStatusQuestionByTrainer implements QuestionUseCase.
// masih belum bisa
func (q *questionUseCase) UpdatadStatusQuestionByTrainer(trainerId string, participantId string, payload dto.QuestionDTO) (dto.QuestionDTO, error) {
	//getSchedule hari H
	logger := logrus.New()
	logger.Infoln(participantId)
	var dates []time.Time
	var day time.Time
	scheduleIDs, err := q.scheduleRepo.GetScheduleByTrainerIdWithoutPagination(trainerId)
	if err != nil {
		return dto.QuestionDTO{}, fmt.Errorf("Field can't be empty")
	}
	for _, scheduleId := range scheduleIDs {
		date, err := q.scheduleRepo.GetDateByScheduleId(scheduleId)
		if err != nil {
			fmt.Println("error tidak dapat menemukan tanggal")
			return dto.QuestionDTO{}, fmt.Errorf("Field can't be empty")
		}
		dates = append(dates, date)
	}
	for _, scheduleDate := range dates {
		//ngetes aja, harusnya pakai time.now biar akurat cuma belum harinya aja
		if scheduleDate.Format("2006-01-02") != time.Date(2023, time.December, 26, 0, 0, 0, 0, time.UTC).Format("2006-01-02") {
			fmt.Println("Anda tidak memiliki jadwal hari ini")
		} else {
			day = scheduleDate
			fmt.Println("Anda memiliki jadwal hari ini")
		}
	}
	//validasi input payload
	if payload.Answer == "" || payload.Status == "" {
		return dto.QuestionDTO{}, fmt.Errorf("oops, Required field is empty")
	}

	schedule, _ := q.scheduleRepo.GetScheduleIdByDate(day)
	payload.ScheduleID = schedule
	payload.TrainerID = trainerId

	//ngecek apakah nama yang bertanya sudah berstatus answered atau belum
	questionData, err := q.repo.GetQuestionByScheduleIdandParticipantId(schedule, participantId)
	if err != nil {
		return dto.QuestionDTO{}, fmt.Errorf("gagal mengambil data")
	}
	if questionData.Answer != "" && questionData.Status != "" {
		return dto.QuestionDTO{}, nil
	}
	logger.Infoln(questionData.ParticipantID)
	logger.Infoln(questionData.ID)
	logger.Infoln(questionData.ScheduleID)
	data, err := q.repo.UpdateQuestionStatusByTrainer(participantId, payload)
	if err != nil {
		return dto.QuestionDTO{}, fmt.Errorf("gagal update")
	}
	return data, nil
}

// CreateQuestionByTrainer implements QuestionUseCase.
func (q *questionUseCase) CreateQuestionByTrainer(trainerId string, participantId string, payload dto.QuestionDTO) (dto.QuestionDTO, error) {
	//getSchedule hari H
	var dates []time.Time
	var day time.Time
	scheduleIDs, err := q.scheduleRepo.GetScheduleByTrainerIdWithoutPagination(trainerId)
	if err != nil {
		return dto.QuestionDTO{}, fmt.Errorf("Field can't be empty")
	}
	fmt.Println(scheduleIDs)
	for _, scheduleId := range scheduleIDs {
		date, err := q.scheduleRepo.GetDateByScheduleId(scheduleId)
		if err != nil {
			fmt.Println("error tidak dapat menemukan tanggal")
			return dto.QuestionDTO{}, fmt.Errorf("Field can't be empty")
		}
		dates = append(dates, date)
	}
	for _, scheduleDate := range dates {
		//ngetes aja, harusnya pakai time.now biar akurat cuma belum harinya aja
		if scheduleDate.Format("2006-01-02") != time.Date(2023, time.December, 26, 0, 0, 0, 0, time.UTC).Format("2006-01-02") {
			fmt.Println("Anda tidak memiliki jadwal hari ini")
		} else {
			day = scheduleDate
		}
	}
	//validasi input payload
	if payload.Question == "" || payload.Answer == "" || payload.Status == "" {
		return dto.QuestionDTO{}, fmt.Errorf("oops, Required field is empty")
	}
	schedule, _ := q.scheduleRepo.GetScheduleIdByDate(day)
	payload.ScheduleID = schedule
	payload.TrainerID = trainerId
	data, err := q.repo.CreateQuestionByTrainer(participantId, payload)
	return data, nil

}

// CreateQuestionByParticipant implements QuestionUseCase.
func (q *questionUseCase) CreateQuestionByParticipant(participantId string, payload dto.QuestionDto) (dto.QuestionDto, error) {
	var dates []time.Time
	var day time.Time
	scheduleIDs, err := q.scheduleRepo.GetScheduleWithParticipantId(participantId)
	if err != nil {
		return dto.QuestionDto{}, fmt.Errorf("field can't be empty")
	}
	if scheduleIDs == nil {
		return dto.QuestionDto{}, fmt.Errorf("received nil schedule IDs")
	}
	for _, scheduleId := range scheduleIDs {
		date, err := q.scheduleRepo.GetDateByScheduleId(scheduleId.ID)
		if err != nil {
			fmt.Println("error tidak dapat menemukan tanggal")
			return dto.QuestionDto{}, fmt.Errorf("field can't be empty")
		}
		dates = append(dates, date)
	}
	for _, scheduleDate := range dates {

		if scheduleDate.Format("2006-01-02") != time.Date(2023, time.December, 26, 0, 0, 0, 0, time.UTC).Format("2006-01-02") {
			fmt.Println("Anda tidak memiliki jadwal hari ini")
		} else {
			day = scheduleDate
		}
	}
	//validasi input payload
	if payload.Question == "" {
		return dto.QuestionDto{}, fmt.Errorf("oops, Required field is empty")
	}
	schedule, err := q.scheduleRepo.GetScheduleIdByDate(day)
	if err != nil {
		return dto.QuestionDto{}, fmt.Errorf("Failed to get schedule ID for the current date: %v", err)
	}
	payload.ScheduleID = schedule

	// Create question
	data, err := q.repo.CreateQuestionByParticipant(participantId, payload)
	if err != nil {
		return dto.QuestionDto{}, fmt.Errorf("Failed to create question")
	}

	return data, nil
}

// FindQuestionByScheduleId implements QuestionUseCase.
func (q *questionUseCase) FindQuestionByTrainerId(id string, page, size int) ([]entity.Question, model.Paging, error) {
	question, paging, err := q.repo.GetQuestionByTrainerId(id, page, size)

	if err != nil {
		return []entity.Question{}, model.Paging{}, fmt.Errorf("failed to get question's : %s", err.Error())

	}

	if len(question) == 0 {
		return []entity.Question{}, model.Paging{}, fmt.Errorf("no question found for the trainer")
	}

	return question, paging, nil
}

func (q *questionUseCase) FindById(id string) (entity.Question, error) {
	return q.repo.Get(id)
}

func (q *questionUseCase) FindAllQuestion(page, size int) ([]entity.Question, model.Paging, error) {
	return q.repo.List(page, size)
}

func (q *questionUseCase) DeleteQuestion(id string) error {
	return q.repo.Delete(id)
}

func (q *questionUseCase) UpdateQuestion(id string, payload entity.Question) (entity.Question, error) {
	payload.UpdatedAt = time.Now()
	question, err := q.repo.Update(id, payload)
	if err != nil {
		return entity.Question{}, fmt.Errorf("failed to update question : %s", err.Error())
	}
	return question, nil
}

func (q *questionUseCase) CreateNewQuestion(payload entity.Question) (entity.Question, error) {
	_, err := q.participantUC.GetParticipantByID(payload.ParticipantID)
	if err != nil {
		return entity.Question{}, fmt.Errorf("participant with ID %s not found", payload.ParticipantID)
	}
	if payload.Question == "" {
		return entity.Question{}, fmt.Errorf("question can't be empty")
	}
	payload.UpdatedAt = time.Now()
	question, err := q.repo.Create(payload)
	if err != nil {
		return entity.Question{}, fmt.Errorf("failed to create new question: %s", err.Error())

	}
	return question, nil
}

func NewQuestionUseCase(repo repository.QuestionRepository, participantRepo repository.ParticipantRepository, scheduleRepo repository.ScheduleRepository, userRepo repository.UserRepository, trainerRepo repository.TrainerRepository, participantUC ParticipantUseCase, trainerUC TrainerUsecase) QuestionUseCase {
	return &questionUseCase{repo: repo, participantRepo: participantRepo, scheduleRepo: scheduleRepo, userRepo: userRepo, trainerRepo: trainerRepo, participantUC: participantUC, trainerUseCase: trainerUC}
}
