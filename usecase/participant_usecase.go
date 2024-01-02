package usecase

import (
	"fmt"
	"instructor-led-app/entity/dto"
	"instructor-led-app/repository"
	"instructor-led-app/shared/model"
	"reflect"
	"strings"
	"time"
)

type ParticipantUseCase interface {
	CreateNewParticipant(participant dto.ParticipantDTO) (dto.ParticipantDTO, error)
	GetAllParticipants(page, size int) ([]dto.ParticipantDTO, model.Paging, error)
	GetParticipantByID(id string) (dto.ParticipantDTO, error)
	UpdateParticipantByID(participant dto.ParticipantDTO) (dto.ParticipantDTO, error)
	DeleteParticipantByID(id string) error
	GetParticipantByUserId(userId string) (dto.ParticipantDTO, error)
	FindScheduleWithParticipantId(userId string) ([]dto.ParticipantDTO, error)
	UpdateParticipantByRole(role, id string) error
}

type participantUseCase struct {
	participantRepository repository.ParticipantRepository
}

// UpdateParticipantByRole implements ParticipantUseCase.
func (p *participantUseCase) UpdateParticipantByRole(role, id string) error {
	if role == "" {
		return fmt.Errorf("role can't empty")
	}
	return p.participantRepository.UpdateByRole(role, id)
}

// GetParticipantByUserId implements ParticipantUseCase.
func (u *participantUseCase) GetParticipantByUserId(userId string) (dto.ParticipantDTO, error) {
	return u.participantRepository.GetParticipantByUserId(userId)
}

// FindScheduleWithParticipantId implements ParticipantUseCase. fitur puji
func (u *participantUseCase) FindScheduleWithParticipantId(userId string) ([]dto.ParticipantDTO, error) {
	participant, err := u.participantRepository.GetParticipantByUserId(userId)
	fmt.Println(participant)
	if err != nil {
		return []dto.ParticipantDTO{}, fmt.Errorf("failed to get participant: %v", err)
	}

	schedul, err := u.participantRepository.GetScheduleById(participant.ID)
	// fmt.Println(schedul)
	if err != nil {
		return []dto.ParticipantDTO{}, fmt.Errorf("failed to get participant id: %v", err)
	}

	if len(schedul) == 0 {
		return []dto.ParticipantDTO{}, fmt.Errorf("schedul not found: %v", err)
	}
	return schedul, nil
}

func (u *participantUseCase) CreateNewParticipant(participant dto.ParticipantDTO) (dto.ParticipantDTO, error) {
	participant, err := u.participantRepository.Insert(participant)
	if err != nil {
		return dto.ParticipantDTO{}, err
	}

	participant.DateOfBirth, err = u.parsedDate(strings.Split(participant.DateOfBirth, "T")[0])
	if err != nil {
		return dto.ParticipantDTO{}, err
	}

	return participant, nil
}

func (u *participantUseCase) parsedDate(value string) (string, error) {
	parsedDate, err := time.Parse("2006-01-02", value)
	if err != nil {
		return "", err
	}

	return parsedDate.Format("2006-01-02"), nil
}

func (u *participantUseCase) GetAllParticipants(page, size int) ([]dto.ParticipantDTO, model.Paging, error) {
	return u.participantRepository.FindAll(page, size)
}

func (u *participantUseCase) GetParticipantByID(id string) (dto.ParticipantDTO, error) {
	return u.participantRepository.FindByID(id)
}

func (u *participantUseCase) UpdateParticipantByID(participant dto.ParticipantDTO) (dto.ParticipantDTO, error) {
	result, err := u.participantRepository.FindByID(participant.ID)
	if err != nil {
		return dto.ParticipantDTO{}, err
	}

	participantMap := u.structToMap(participant)
	resultMap := u.structToMap(result)

	// Mengisi nilai participantMap dengan nilai resultMap jika ada data yang kosong
	for key, value := range participantMap {
		if *value == "" {
			*value = *resultMap[key]
		}
	}

	// Mengisi nilai participantMap ke struct participant
	for key, value := range participantMap {
		field := reflect.ValueOf(&participant).Elem().FieldByName(key)

		// Check if the field type is string before setting the value
		if field.Kind() == reflect.String {
			fieldValue := *value // Dereference the pointer
			field.SetString(fieldValue)
		}
	}

	participantDto, err := u.participantRepository.UpdateByID(participant, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return dto.ParticipantDTO{}, err
	}

	participantDto.DateOfBirth, _ = u.parsedDate(strings.Split(participantDto.DateOfBirth, "T")[0])

	return participantDto, nil
}

func (u *participantUseCase) structToMap(s dto.ParticipantDTO) map[string]*string {
	structValue := reflect.ValueOf(s)
	structType := structValue.Type()

	result := make(map[string]*string)

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldName := structType.Field(i).Name
		fieldValue := field.String()
		result[fieldName] = &fieldValue
	}

	return result
}

func (u *participantUseCase) DeleteParticipantByID(id string) error {
	return u.participantRepository.DeleteByID(id)
}

func NewParticipantUseCase(participantRepository repository.ParticipantRepository) ParticipantUseCase {
	return &participantUseCase{participantRepository: participantRepository}
}
