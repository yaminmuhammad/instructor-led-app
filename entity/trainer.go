package entity

import "time"

type Trainer struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	//Specializations []string  `json:"specializations"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
