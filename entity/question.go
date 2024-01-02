package entity

import "time"

type Question struct {
	ID            string    `json:"id"`
	Question      string    `json:"question"`
	Answer        string    `json:"answer"`
	Status        string    `json:"status"`
	ParticipantID string    `json:"participantId"`
	TrainerID     string    `json:"trainerId"`
	ScheduleID    string    `json:"scheduleId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
