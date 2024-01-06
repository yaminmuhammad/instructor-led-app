package usecase

import (
	"fmt"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/repository"
	"instructor-led-app/shared/model"
	"log"
	"time"
)

type TrainerUsecase interface {
	FindAllTrainer(page, size int) ([]entity.Trainer, model.Paging, error)
	FindTrainerById(trainerId string) ([]entity.Trainer, error)
	FindTrainerByUserIDs(userID string) (dto.TrainerDTO, error)
	TrainerUpdated(trainer dto.TrainerDTO) (dto.TrainerDTO, error)
	DeleteTrainer(trainerId string) (entity.Trainer, error)
	FindTrainerByUserId(userId string) (entity.Trainer, error)
}

type trainerUseCase struct {
	repo repository.TrainerRepository
}

// FindTrainerByUserId implements TrainerUsecase.
func (t *trainerUseCase) FindTrainerByUserId(userId string) (entity.Trainer, error) {
	return t.repo.TrainerByUserId(userId)
}

// TrainerUpdated implements TrainerUsecase.
func (t *trainerUseCase) TrainerUpdated(trainer dto.TrainerDTO) (dto.TrainerDTO, error) {
	if trainer.ID == "" && trainer.PhoneNumber == "" && trainer.UserID == "" {
		return dto.TrainerDTO{}, fmt.Errorf("field can't be empty")
	}
	return t.repo.UpdateTrainer(trainer, time.Now())
}

// Delete implements TrainerUsecase.
func (t *trainerUseCase) DeleteTrainer(trainerId string) (entity.Trainer, error) {
	trainer, err := t.repo.Delete(trainerId)
	if err != nil {
		log.Fatal("delete error")
		return entity.Trainer{}, err
	}
	return trainer, nil
}

func (t *trainerUseCase) FindTrainerById(trainerId string) ([]entity.Trainer, error) {
	return t.repo.TrainerById(trainerId)
}

func (t *trainerUseCase) FindTrainerByUserIDs(userID string) (dto.TrainerDTO, error) {
	return t.repo.FindByUserID(userID)
}

// FindAllTrainer implements TrainerUsecase.
func (t *trainerUseCase) FindAllTrainer(page, size int) ([]entity.Trainer, model.Paging, error) {
	return t.repo.List(page, size)
}

func NewTrainerUseCase(repo repository.TrainerRepository) TrainerUsecase {
	return &trainerUseCase{repo: repo}
}
