package dto

import "time"

type AbsenceDTO struct {
	ID             string    `json:"id"`
	Date           time.Time `json:"date"`
	Schedule_id    string    `json:"scheduleId"`
	Trainer_id     string    `json:"trainerId"`
	Participant_id string    `json:"participantId"`
}
