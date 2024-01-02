package entity

import "time"

type User struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Address      string    `json:"address"`
	Hashpassword string    `json:"hashPassword"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (e User) IsroleTypeValid() bool {
	return e.Role == "admin" || e.Role == "trainer" || e.Role == "participant"
}
