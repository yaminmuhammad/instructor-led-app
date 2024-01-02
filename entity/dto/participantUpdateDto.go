package dto

type ParticipantUpdateDTO struct {
	DateOfBirth   string `json:"dateOfBirth"`
	PlaceOfBirth  string `json:"placeOfBirth"`
	LastEducation string `json:"lastEducation"`
}
