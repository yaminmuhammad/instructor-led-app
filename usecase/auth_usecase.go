package usecase

import (
	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/service"
)

type AuthUseCase interface {
	Login(payload dto.AuthRequestDto) (dto.AuthResponseDto, error)
}

type authUseCase struct {
	userUC     UserUsecase
	jwtService service.JwtService
}

func (a *authUseCase) Login(payload dto.AuthRequestDto) (dto.AuthResponseDto, error) {
	user, err := a.userUC.AuthUser(payload.Email, payload.HashPassword)
	if err != nil {
		return dto.AuthResponseDto{}, err
	}
	token, err := a.jwtService.CreateToken(user)
	if err != nil {
		return dto.AuthResponseDto{}, err
	}

	return token, nil
}

func NewAuthUseCase(userUC UserUsecase, jwtService service.JwtService) AuthUseCase {
	return &authUseCase{userUC: userUC, jwtService: jwtService}
}
