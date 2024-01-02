package dto

import "time"

type ParticipantScheduleDTO struct {
	ID         string      `json:"id"`
	TrainerID  string      `json:"trainerId"`
	Date       []time.Time `json:"Date"`
	ScheduleID []string    `json:"scheduleId"`
}
