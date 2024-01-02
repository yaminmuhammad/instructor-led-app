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

type AbsenceController struct {
	absenceUC      usecase.AbsenceUseCase
	scheduleUC     usecase.ScheduleUseCase
	trainerUC      usecase.TrainerUsecase
	participantUC  usecase.ParticipantUseCase
	userUC         usecase.UserUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (a *AbsenceController) createHandler(ctx *gin.Context) {
	var payload dto.TrainerNameDTO
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	absence, err := a.absenceUC.InsertNewAbsence(payload.Name)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	common.SendCreateResponse(ctx, absence, "Created")
}

func (a *AbsenceController) UpdateAbsencesByParticipantName(ctx *gin.Context) {
	var payload dto.AbsenceCheckDTO
	user := ctx.MustGet("userID").(string)
	trainer, _ := a.trainerUC.FindTrainerByUserId(user)
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, _ := a.userUC.FindUserIDByName(payload.Name)
	participantId, _ := a.participantUC.GetParticipantByUserId(userId.Id)

	updateAbsence, err := a.absenceUC.UpdateAbsencesByScheduleId(trainer.ID, participantId.ID, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed update"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": updateAbsence})

}

func (a *AbsenceController) GetAbsencesByScheduleIdHandler(ctx *gin.Context) {
	user := ctx.MustGet("userID").(string)
	trainer, err := a.trainerUC.FindTrainerByUserId(user)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusNotFound, "Trainer tidak ada")
		return
	}
	absences, err := a.absenceUC.GetAbsencesByScheduleID(trainer.ID)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusNotFound, "Anda tidak memiliki jadwal hari ini")
		return
	}
	common.SendSingleResponse(ctx, absences, "Ok")

}

func (a *AbsenceController) listHandler(ctx *gin.Context) {
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

	absences, paging, err := a.absenceUC.FindAllAbsence(startDate, endDate, page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var response []interface{}
	for _, v := range absences {
		response = append(response, v)
	}
	common.SendPagedResponse(ctx, response, paging, "Ok")
}

func (a *AbsenceController) GetByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	absences, err := a.absenceUC.GetAbsencesByParticipantID(id)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusNotFound, "task with author ID "+id+" not found")
		return
	}
	common.SendSingleResponse(ctx, absences, "Ok")
}

func (a *AbsenceController) deleteHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	// Call the DeleteByParticipantId method from the use case
	err := a.absenceUC.DeleteByParticipantId(id)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with a success message
	common.SendDeleteResponse(ctx, "Absence deleted successfully")
}

func (a *AbsenceController) Route() {
	a.rg.POST(config.AbsencePost, a.authMiddleware.RequireToken("admin"), a.createHandler) //dah bisa tapi logicnya baru masuk akal kalau data participant id di schedule diilangin
	a.rg.GET(config.ListAbsences, a.authMiddleware.RequireToken("admin", "trainer"), a.listHandler)
	a.rg.GET(config.ListAbsencesById, a.authMiddleware.RequireToken("admin", "trainer"), a.GetByIdHandler)
	a.rg.DELETE(config.DeleteAbsence, a.authMiddleware.RequireToken("admin"), a.deleteHandler)
	a.rg.GET(config.AbsenceByTrainerScheduleId, a.authMiddleware.RequireToken("trainer"), a.GetAbsencesByScheduleIdHandler)
	a.rg.PUT(config.AbsenceParticipantByTrainer, a.authMiddleware.RequireToken("trainer"), a.UpdateAbsencesByParticipantName)
}

func NewAbsenceController(absenceUc usecase.AbsenceUseCase, scheduleUc usecase.ScheduleUseCase, trainerUc usecase.TrainerUsecase, participantUC usecase.ParticipantUseCase, userUC usecase.UserUsecase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *AbsenceController {
	return &AbsenceController{
		absenceUC:      absenceUc,
		scheduleUC:     scheduleUc,
		trainerUC:      trainerUc,
		participantUC:  participantUC,
		userUC:         userUC,
		rg:             rg,
		authMiddleware: authMiddleware,
	}
}
