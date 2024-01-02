package usecase

import (
	"fmt"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/repository"
	"instructor-led-app/shared/model"
	"time"
)

type AbsenceUseCase interface {
	InsertNewAbsence(name string) ([]dto.ParticipantScheduleDTO, error)
	FindAllAbsence(startDate, endDate time.Time, page, size int) ([]entity.Absence, model.Paging, error)
	GetAbsencesByParticipantID(id string) (entity.Absence, error)
	GetAbsencesByScheduleID(id string) (dto.AbsenceDTO, error)
	UpdateAbsencesByScheduleId(trainerId, participantId string, payload dto.AbsenceCheckDTO) (dto.AbsenceCheckDTO, error)
	DeleteByParticipantId(id string) error
}

type absenceUseCase struct {
	repo            repository.AbsenceRepository
	participantRepo repository.ParticipantRepository
	scheduleRepo    repository.ScheduleRepository
	userRepo        repository.UserRepository
	trainerRepo     repository.TrainerRepository
}

// UpdateAbsencesByScheduleId implements AbsenceUseCase.
func (a *absenceUseCase) UpdateAbsencesByScheduleId(trainerId, participantId string, payload dto.AbsenceCheckDTO) (dto.AbsenceCheckDTO, error) {
	var dates []time.Time
	var day time.Time
	scheduleIDs, err := a.scheduleRepo.GetScheduleByTrainerIdWithoutPagination(trainerId)
	fmt.Print(trainerId)
	fmt.Print(scheduleIDs)
	if err != nil {
		return dto.AbsenceCheckDTO{}, fmt.Errorf("Field can't be empty")
	}
	fmt.Println(scheduleIDs)
	for _, scheduleId := range scheduleIDs {
		date, err := a.scheduleRepo.GetDateByScheduleId(scheduleId)
		if err != nil {
			fmt.Println("error tidak dapat menemukan tanggal")
			return dto.AbsenceCheckDTO{}, fmt.Errorf("Field can't be empty")
		}
		dates = append(dates, date)
	}
	for _, scheduleDate := range dates {
		//ngetes aja, harusnya pakai time.now biar akurat cuma belum harinya aja
		if scheduleDate.Format("2006-01-02") != time.Date(2023, time.December, 19, 0, 0, 0, 0, time.UTC).Format("2006-01-02") {
			fmt.Println("Anda tidak memiliki jadwal hari ini")
		} else {
			day = scheduleDate
		}
	}
	absenceTime := time.Now().Format("2006-01-2 15:04:03")
	absenceTimeReal, _ := time.Parse("2006-01-2 15:04:05", absenceTime)
	updatedAtTime := time.Now().Format("2006-01-2 15:04:03")
	updatedAtTimeReal, _ := time.Parse("2006-01-2 15:04:05", updatedAtTime)
	payload.Updated_at = updatedAtTimeReal
	payload.Absence_time = absenceTimeReal
	schedule, _ := a.scheduleRepo.GetScheduleIdByDate(day)
	data, err := a.repo.UpdateByScheduleIDandParticipantID(schedule, participantId, payload)
	return data, nil
}

// GetAbsencesByScheduleID implements AbsenceUseCase.
func (a *absenceUseCase) GetAbsencesByScheduleID(id string) (dto.AbsenceDTO, error) {
	var dates []time.Time
	var day time.Time
	scheduleIDs, err := a.scheduleRepo.GetScheduleByTrainerIdWithoutPagination(id)
	if err != nil {
		return dto.AbsenceDTO{}, err
	}
	for _, scheduleId := range scheduleIDs {
		date, err := a.scheduleRepo.GetDateByScheduleId(scheduleId)
		if err != nil {
			fmt.Println("error tidak dapat menemukan tanggal")
			return dto.AbsenceDTO{}, err
		}
		dates = append(dates, date)
	}
	for _, scheduleDate := range dates {
		//ngetes aja, harusnya pakai time.now biar akurat cuma belum harinya aja
		if scheduleDate.Format("2006-01-02") != time.Date(2023, time.December, 19, 0, 0, 0, 0, time.UTC).Format("2006-01-02") {
			fmt.Println("Anda tidak memiliki jadwal hari ini")
		} else {
			day = scheduleDate
		}
	}

	schedule, _ := a.scheduleRepo.GetScheduleIdByDate(day)
	absences, err := a.repo.GetAbsencesByScheduleID(schedule)
	if err != nil {
		return dto.AbsenceDTO{}, err
	}
	return absences, err
}

// DeleteByParticipantId implements AbsenceUseCase.
func (a *absenceUseCase) DeleteByParticipantId(id string) error {
	err := a.repo.Delete(id)
	if err != nil {
		fmt.Println("error deleting absence:", err)
		return err
	}
	return nil
}

// GetAbsencesByParticipantID implements AbsenceUseCase.
func (a *absenceUseCase) GetAbsencesByParticipantID(id string) (entity.Absence, error) {
	return a.repo.GetAbsencesByParticipantID(id)
}

// FindAllTask implements AbsenceUseCase.
func (a *absenceUseCase) FindAllAbsence(startDate time.Time, endDate time.Time, page int, size int) ([]entity.Absence, model.Paging, error) {
	if !startDate.IsZero() && !endDate.IsZero() {
		return a.repo.ListByDate(startDate, endDate, page, size)
	} else {
		return a.repo.List(page, size)
	}
}

// InsertNewAbsence implements AbsenceUseCase.
func (a *absenceUseCase) InsertNewAbsence(name string) ([]dto.ParticipantScheduleDTO, error) {
	trainer, err := a.userRepo.GetUserByName(name)
	if err != nil {
		// Handle error
		fmt.Println("Error fetching userId", err)
		return nil, err
	}
	trainerUserId := trainer.Id
	trainers, err := a.trainerRepo.TrainerByUserId(trainerUserId)
	trainerID := trainers.ID
	if err != nil {
		// Handle error
		fmt.Println("Error fetching trainerId:", err)
		return nil, err
	}
	var userScheduleDto []dto.ParticipantScheduleDTO
	participants, err := a.participantRepo.FindAllIDandRole()
	if err != nil {
		// Handle error
		fmt.Println("Error fetching participants:", err)
		return nil, fmt.Errorf("oppps, failed to save data absence :%v", err.Error())
	}

	roleToDay := map[string]int{
		"Basic":   1, // Monday
		"Advance": 2, // Tuesday
		// Add more roles and corresponding days as needed
	}

	for _, participant := range participants {
		day, found := roleToDay[participant.Role]

		if !found {
			fmt.Printf("Unknown role: %s\n", participant.Role)
			continue
		}
		// Get the schedule ID using the scheduleRepo
		scheduleIds, _ := a.scheduleRepo.GetScheduleIdByDay(day)
		date, _ := a.scheduleRepo.GetDateByDay(day)

		userSchedule := dto.ParticipantScheduleDTO{
			ID:         participant.ID,
			TrainerID:  trainerID,
			Date:       date,
			ScheduleID: scheduleIds,
		}
		userScheduleDto = append(userScheduleDto, userSchedule)
	}

	absence, err := a.repo.Create(userScheduleDto)
	if err != nil {
		return nil, fmt.Errorf("oppps, failed to save data absence :%v", err.Error())
	}
	return absence, nil
}

func NewAbsenceUseCase(repo repository.AbsenceRepository, participantRepo repository.ParticipantRepository, scheduleRepo repository.ScheduleRepository, userRepo repository.UserRepository, trainerRepo repository.TrainerRepository) AbsenceUseCase {
	return &absenceUseCase{repo: repo, participantRepo: participantRepo, scheduleRepo: scheduleRepo, userRepo: userRepo, trainerRepo: trainerRepo}
}
