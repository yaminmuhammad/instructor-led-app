package middleware

import (
	"instructor-led-app/repository"
	"instructor-led-app/shared/service"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
	repo       repository.ParticipantRepository
}

type AuthHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var autHeader AuthHeader
		if err := ctx.ShouldBindHeader(&autHeader); err != nil {
			log.Printf("RequireToken.autHeader: %v \n", err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenHeader := strings.Replace(autHeader.AuthorizationHeader, "Bearer ", "", -1)
		if tokenHeader == "" {
			log.Printf("RequireToken.tokenHeader \n")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := a.jwtService.ParseToken(tokenHeader)
		if err != nil {
			log.Printf("RequireToken.ParseToken: %v \n", err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("userID", claims["userId"])

		validRole := false
		// admin, user, other....
		for _, role := range roles {
			if role == claims["role"] {
				validRole = true
				break
			}
		}

		if !validRole {
			log.Printf("RequireToken.validRole\n")
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
