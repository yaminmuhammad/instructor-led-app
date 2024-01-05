package delivery

import (
	"database/sql"
	"fmt"
	"instructor-led-app/config"
	"instructor-led-app/delivery/controller"
	"instructor-led-app/delivery/middleware"
	"instructor-led-app/repository"
	"instructor-led-app/shared/service"
	"instructor-led-app/usecase"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	trainerUseCase       usecase.TrainerUsecase
	participantUseCase   usecase.ParticipantUseCase
	userUc               usecase.UserUsecase
	authUc               usecase.AuthUseCase
	questionUc           usecase.QuestionUseCase
	absenceUC            usecase.AbsenceUseCase
	scheduleUC           usecase.ScheduleUseCase
	scheduleImageUseCase usecase.ScheduleImageUseCase
	jwtService           service.JwtService
	engine               *gin.Engine
	port                 string
}

func (s *Server) initRoute() {
	rg := s.engine.Group(config.APIGroup)

	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)
	controller.NewTrainerController(s.trainerUseCase, rg, authMiddleware).Route()
	controller.NewParticipantController(s.participantUseCase, s.userUc, rg, authMiddleware).Route()
	controller.NewAuthController(s.authUc, rg).Route()
	controller.NewUserController(s.userUc, rg, authMiddleware).Route()
	controller.NewAbsenceController(s.absenceUC, s.scheduleUC, s.trainerUseCase, s.participantUseCase, s.userUc, rg, authMiddleware).Route()
	controller.NewQuestionController(s.questionUc, s.scheduleUC, s.trainerUseCase, s.participantUseCase, s.userUc, rg, authMiddleware).Route()
	controller.NewScheduleController(s.scheduleUC, s.userUc, s.trainerUseCase, rg, authMiddleware).Route()
	controller.NewScheduleImageController(s.scheduleImageUseCase, authMiddleware, rg).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.port); err != nil {
		log.Fatalf("server can't running on port '%v', error : %v", s.port, err)
	}
}

func NewServer() *Server {
	config, _ := config.NewConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Name)
	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		panic("connection error")
	}

	//repo
	userRepo := repository.NewUserRepository(db)
	absenceRepo := repository.NewAbsenceRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	trainerRepo := repository.NewTrainerRepository(db)

	jwtService := service.NewJwtService(config.TokenConfig)
	participantRepository := repository.NewParticipantRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	scheduleImageRepository := repository.NewScheduleImagesRepository(db)
	// usecase
	trainerUseCase := usecase.NewTrainerUseCase(trainerRepo)
	absenceUC := usecase.NewAbsenceUseCase(absenceRepo, participantRepository, scheduleRepo, userRepo, trainerRepo)
	participantUseCase := usecase.NewParticipantUseCase(participantRepository)
	UserUsecase := usecase.NewUserUsecase(userRepo)
	questionUsecase := usecase.NewQuestionUseCase(questionRepo, participantRepository, scheduleRepo, userRepo, trainerRepo, participantUseCase, trainerUseCase)
	scheduleUC := usecase.NewScheduleUseCase(scheduleRepo, trainerUseCase)
	scheduleImageUseCase := usecase.NewScheduleImageUseCase(scheduleImageRepository, scheduleUC)

	authUc := usecase.NewAuthUseCase(UserUsecase, jwtService)

	engine := gin.Default()
	port := fmt.Sprintf(":%s", config.ApiPort)
	return &Server{
		trainerUseCase,
		participantUseCase,
		UserUsecase,
		authUc,
		questionUsecase,
		absenceUC,
		scheduleUC,
		scheduleImageUseCase,
		jwtService,
		engine,
		port,
	}
}
