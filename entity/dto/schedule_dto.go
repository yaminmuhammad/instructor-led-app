package dto

import "time"

type ScheduleDto struct {
	ID            string        `json:"id"`
	Activity      string        `json:"activity"`
	Date          string        `json:"date"`
	TrainerID     string        `json:"trainerId"`
	ParticipantID string        `json:"participantId"`
	Day           string        `json:"day"`
	Question      []QuestionDTO `json:"question"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}
