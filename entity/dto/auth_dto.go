package dto

type AuthRequestDto struct {
	Email        string `json:"email"`
	HashPassword string `json:"hashPassword"`
}

type AuthResponseDto struct {
	Token string `json:"token"`
}
