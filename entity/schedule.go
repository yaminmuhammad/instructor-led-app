package entity

import "time"

type Schedule struct {
	ID            string    `json:"id"`
	Activity      string    `json:"activity"`
	Date          time.Time `json:"date"`
	TrainerID     string    `json:"trainerId"`
	ParticipantID string    `json:"participantId"`
	Day           string    `json:"day"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
