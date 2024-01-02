package controller

import (
	"instructor-led-app/config"
	"instructor-led-app/delivery/middleware"
	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/common"
	"instructor-led-app/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var ScheduleId []string

type ScheduleController struct {
	scheduleUC     usecase.ScheduleUseCase
	userUC         usecase.UserUsecase
	trainerUC      usecase.TrainerUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (s *ScheduleController) UpdateByAdminHandler(ctx *gin.Context) {
	var payload dto.UpdateAdminDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := s.userUC.FindUserIDByName(payload.TrainerName)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "No userId Found")
		return
	}
	trainerId, err := s.trainerUC.FindTrainerByUserId(userId.Id)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "No trainer Found")
		return
	}
	updateSchedule, err := s.scheduleUC.UpdateScheduleByAdmin(trainerId.ID, payload.CodeDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed update"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": updateSchedule})

}

func (s *ScheduleController) createScheduleHandler(ctx *gin.Context) {
	var payload dto.ScheduleDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	schedule, err := s.scheduleUC.InsertNewSchedule(payload)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	common.SendCreateResponse(ctx, schedule, "Created")
}

func (s *ScheduleController) listScheduleHandler(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))

	// Tambahkan variabel untuk startDate dan endDate
	startDateParam := ctx.Query("startDate")
	endDateParam := ctx.Query("endDate")

	var startDate, endDate time.Time
	var err error

	// Parsing startDate
	if startDateParam != "" {
		startDate, err = time.Parse("2006-01-02", startDateParam)
		if err != nil {
			common.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid startDate format")
			return
		}
	}

	// Parsing endDate
	if endDateParam != "" {
		endDate, err = time.Parse("2006-01-02", endDateParam)
		if err != nil {
			common.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid endDate format")
			return
		}
	}
	schedules, paging, err := s.scheduleUC.FindAllSchedule(startDate, endDate, page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var response []interface{}
	for _, v := range schedules {
		response = append(response, v)
	}
	common.SendPagedResponse(ctx, response, paging, "Ok")
}

func (s *ScheduleController) GetScheduleByParticipantID(ctx *gin.Context) {
	id := ctx.Param("id")
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	schedules, paging, err := s.scheduleUC.GetScheduleByParticipantID(id, page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var response []interface{}
	for _, v := range schedules {
		response = append(response, v)
	}
	common.SendPagedResponse(ctx, response, paging, "Ok")
}

func (s *ScheduleController) GetScheduleByTrainerID(ctx *gin.Context) {
	userId := ctx.MustGet("userID").(string)
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	schedules, paging, err := s.scheduleUC.GetScheduleByTrainerID(userId, page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var response []interface{}
	for _, v := range schedules {
		response = append(response, v)
	}
	common.SendPagedResponse(ctx, response, paging, "Ok")
}

func (s *ScheduleController) deleteHandler(ctx *gin.Context) {
	date := ctx.Param("date")

	if err := s.scheduleUC.DeleteScheduleByDate(date); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	common.SendDeleteResponse(ctx, "Delete Schedule successfully")
}

func (s *ScheduleController) Route() {
	s.rg.POST(config.SchedulePost, s.authMiddleware.RequireToken("admin"), s.createScheduleHandler)           //taran
	s.rg.GET(config.ScheduleList, s.authMiddleware.RequireToken("admin"), s.listScheduleHandler)              //bisa
	s.rg.GET(config.ScheduleById, s.authMiddleware.RequireToken("participant"), s.GetScheduleByParticipantID) //ini harusnya bagian puji sih
	s.rg.GET(config.ScheduleByTrainerId, s.authMiddleware.RequireToken("trainer"), s.GetScheduleByTrainerID)  //bisa
	s.rg.PUT(config.ScheduleByTrainerId, s.authMiddleware.RequireToken("admin"), s.UpdateByAdminHandler)      //bisa
	s.rg.DELETE(config.DeleteSchedule, s.authMiddleware.RequireToken("admin"), s.deleteHandler)
}

func NewScheduleController(scheduleUc usecase.ScheduleUseCase, userUC usecase.UserUsecase, trainerUC usecase.TrainerUsecase, rg *gin.RouterGroup, auth middleware.AuthMiddleware) *ScheduleController {
	return &ScheduleController{
		scheduleUC:     scheduleUc,
		userUC:         userUC,
		trainerUC:      trainerUC,
		rg:             rg,
		authMiddleware: auth,
	}
}
