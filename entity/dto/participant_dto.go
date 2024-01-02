package dto

type ParticipantDTO struct {
	ID            string        `json:"id"`
	DateOfBirth   string        `json:"dateOfBirth"`
	PlaceOfBirth  string        `json:"placeOfBirth"`
	LastEducation string        `json:"lastEducation"`
	UserID        string        `json:"userId"`
	Role          string        `json:"role"`
	Schedules     []ScheduleDto `json:"schedule"`
}

type ParticipantTestDTO struct {
	ID            string `json:"id"`
	DateOfBirth   string `json:"dateOfBirth"`
	PlaceOfBirth  string `json:"placeOfBirth"`
	LastEducation string `json:"lastEducation"`
	UserID        string `json:"userId"`
	Role          string `json:"role"`
}
