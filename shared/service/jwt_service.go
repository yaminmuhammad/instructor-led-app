package service

import (
	"fmt"
	"instructor-led-app/config"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateToken(user entity.User) (dto.AuthResponseDto, error)
	ParseToken(tokenHeader string) (jwt.MapClaims, error)
}

type jwtService struct {
	cfg config.TokenConfig
}

func (j *jwtService) CreateToken(user entity.User) (dto.AuthResponseDto, error) {
	claims := model.MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.IssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.JwtExpiresTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: user.Id,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
	ss, err := token.SignedString(j.cfg.JwtSignatureKey)
	if err != nil {
		return dto.AuthResponseDto{}, fmt.Errorf("oops, failed to create token")
	}
	return dto.AuthResponseDto{Token: ss}, nil
}

func (j *jwtService) ParseToken(tokenHeader string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return j.cfg.JwtSignatureKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("oops, failed to verify token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("oops, failed to claim token")
	}
	return claims, nil
}

func NewJwtService(cfg config.TokenConfig) JwtService {
	return &jwtService{cfg: cfg}
}
