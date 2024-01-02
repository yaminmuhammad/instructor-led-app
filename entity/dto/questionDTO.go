package dto

import "time"

type QuestionDTO struct {
	ID              string    `json:"id"`
	ParticipantName string    `json:"participantName"`
	Question        string    `json:"question"`
	Answer          string    `json:"answer"`
	Status          string    `json:"status"`
	TrainerID       string    `json:"trainerId"`
	ParticipantID   string    `json:"participantId"`
	ScheduleID      string    `json:"ScheduleId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
type QuestionDto struct {
	ID            string    `json:"id"`
	Question      string    `json:"question"`
	TrainerID     string    `json:"trainerId"`
	ParticipantID string    `json:"participantId"`
	ScheduleID    string    `json:"ScheduleId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
