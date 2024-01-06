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

type TrainerRepository interface {
	List(page, size int) ([]entity.Trainer, model.Paging, error)
	TrainerById(trainerId string) ([]entity.Trainer, error)
	FindByUserID(userID string) (dto.TrainerDTO, error)
	UpdateTrainer(trainerDTO dto.TrainerDTO, updateAt time.Time) (dto.TrainerDTO, error)
	Delete(trainerId string) (entity.Trainer, error)
	TrainerByUserId(userId string) (entity.Trainer, error)
}

type trainerRepository struct {
	db *sql.DB
}

// TrainerByUserId implements TrainerRepository.
func (t *trainerRepository) TrainerByUserId(userId string) (entity.Trainer, error) {
	var trainer entity.Trainer
	err := t.db.QueryRow(config.GetTrainerByUserID, userId).Scan(
		&trainer.ID,
		&trainer.PhoneNumber,
		&trainer.UserID,
		&trainer.CreatedAt, &trainer.UpdatedAt,
	)
	if err != nil {
		return entity.Trainer{}, err
	}
	return trainer, nil
}

// UpdateTrainer implements TrainerRepository.
func (t *trainerRepository) UpdateTrainer(trainerDTO dto.TrainerDTO, updateAt time.Time) (dto.TrainerDTO, error) {
	var trainer dto.TrainerDTO

	if err := t.db.QueryRow(config.UpdateTrainerByID, trainerDTO.PhoneNumber, trainerDTO.UserID, updateAt, trainerDTO.ID).Scan(&trainer.ID, &trainer.PhoneNumber, &trainer.UserID); err != nil {
		log.Println("QueryRow. err Update :", err)
		return dto.TrainerDTO{}, err
	}
	return trainer, nil
}

// Delete implements TrainerRepository.
func (t *trainerRepository) Delete(trainerId string) (entity.Trainer, error) {
	_, err := t.db.Exec(config.DeleteTrainerByID, trainerId)
	if err != nil {
		fmt.Println("error deleting absence:", err)
		return entity.Trainer{}, err
	}
	return entity.Trainer{}, err

}

// TrainerById implements TrainerRepository.
func (t *trainerRepository) TrainerById(trainerId string) ([]entity.Trainer, error) {
	var trainers []entity.Trainer
	rows, err := t.db.Query(config.GetTrainerByID, trainerId)
	if err != nil {
		log.Println("trainerREpository:", err.Error())
		return nil, err
	}
	for rows.Next() {
		var trainer entity.Trainer
		err := rows.Scan(
			&trainer.ID,
			&trainer.PhoneNumber,
			&trainer.UserID,
			&trainer.CreatedAt,
			&trainer.UpdatedAt,
		)
		if err != nil {
			log.Println("trainerRepositoery.GetById.Rows:", err.Error())
			return nil, err
		}
		trainers = append(trainers, trainer)
	}
	return trainers, nil
}

func (t *trainerRepository) FindByUserID(userID string) (dto.TrainerDTO, error) {
	var trainer dto.TrainerDTO
	if err := t.db.QueryRow(config.GetTrainerByUserID, userID).Scan(&trainer.ID, &trainer.PhoneNumber, &trainer.UserID); err != nil {
		log.Println("trainerRepository:", err.Error())
		return dto.TrainerDTO{}, err
	}

	return trainer, nil
}

// List implements TrainerRepository.
func (t *trainerRepository) List(page, size int) ([]entity.Trainer, model.Paging, error) {
	var trainers []entity.Trainer
	offset := (page - 1) * size
	rows, err := t.db.Query(config.ListTrainers, size, offset)

	if err != nil {
		log.Println("TrainerRepository.Query:", err.Error())
		return nil, model.Paging{}, err
	}

	for rows.Next() {
		var trainer entity.Trainer
		err := rows.Scan(
			&trainer.ID,
			&trainer.PhoneNumber,
			&trainer.UserID,
			&trainer.CreatedAt,
			&trainer.UpdatedAt,
		)
		if err != nil {
			log.Println("UserRepository.Rows.Next():", err.Error())
			return nil, model.Paging{}, err
		}

		trainers = append(trainers, trainer)
	}
	totalRows := 0
	if err := t.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalRows); err != nil {
		return nil, model.Paging{}, err
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}
	return trainers, paging, nil

}

// CreateTrainer implements TrainerRepository.

func NewTrainerRepository(db *sql.DB) TrainerRepository {
	return &trainerRepository{db: db}
}
