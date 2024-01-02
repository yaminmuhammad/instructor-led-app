package dto

import "time"

type AbsenceCheckDTO struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Information    string    `json:"information"`
	Absence_status string    `json:"absenceStatus"`
	Absence_time   time.Time `json:"absenceTime"`
	Updated_at     time.Time `json:"updatedAt"`
}
