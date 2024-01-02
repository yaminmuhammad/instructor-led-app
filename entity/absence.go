package entity

import "time"

type Absence struct {
	ID             string    `json:"id"`
	Date           time.Time `json:"date"`
	Information    string    `json:"information"`
	Absence_status string    `json:"absenceStatus"`
	Absence_time   time.Time `json:"absenceTime"`
	Schedule_id    string    `json:"scheduleId"`
	Trainer_id     string    `json:"trainerId"`
	Participant_id string    `json:"participantId"`
	Created_at     time.Time `json:"createdAt"`
	Updated_at     time.Time `json:"updatedAt"`
}
