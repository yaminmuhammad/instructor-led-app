package controller

import (
	"instructor-led-app/config"
	"instructor-led-app/delivery/middleware"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/common"
	"instructor-led-app/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TrainerController struct {
	trainerUc      usecase.TrainerUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (t *TrainerController) createHandle(ctx *gin.Context) {
	var payload entity.Trainer
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	trainer, err := t.trainerUc.RegisterNewTrainer(payload)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	common.SendSingleResponse(ctx, trainer, "Created")
}
func (t *TrainerController) listHandler(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))

	trainer, paging, err := t.trainerUc.FindAllTrainer(page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	var response []interface{}
	for _, v := range trainer {
		response = append(response, v)
	}

	common.SendPagedResponse(ctx, response, paging, "Ok")
}

func (t *TrainerController) trainerByIdHandler(ctx *gin.Context) {
	tranerId := ctx.Param("id")
	if tranerId == "" {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid ID")
		return
	}
	trainer, err := t.trainerUc.FindTrainerById(tranerId)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, "trainer with ID "+tranerId+" not found")
		return
	}
	common.SendSingleResponse(ctx, trainer, "Ok")
}

func (t *TrainerController) trainerByUserIdHandler(ctx *gin.Context) {
	tranerId := ctx.Param("id")
	trainer, err := t.trainerUc.FindTrainerByUserId(tranerId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": trainer})
}

func (t *TrainerController) deleteHandler(ctx *gin.Context) {
	tranerId := ctx.Param("id")

	if tranerId == "" {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid ID")
		return
	}

	_, err := t.trainerUc.DeleteTrainer(tranerId)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete trainer: "+err.Error())
		return
	}
	common.SendSingleResponse(ctx, "Ok", "Delete succesfully")
}

func (t *TrainerController) updatedTrainerHandler(ctx *gin.Context) {
	trainerId := ctx.Param("id")

	if trainerId == "" {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid ID")
		return
	}

	var trainer dto.TrainerDTO

	if err := ctx.ShouldBindJSON(&trainer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		return
	}

	trainer.ID = trainerId

	updateTrainer, err := t.trainerUc.TrainerUpdated(trainer)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to update trainer: "+err.Error())
		return
	}
	common.SendSingleResponse(ctx, updateTrainer, "Ok")
}

func (t *TrainerController) Route() {
	admin := t.rg.Group(config.AdminGroup)

	admin.GET(config.MasterDataTrainerByUserID, t.authMiddleware.RequireToken("admin"), t.trainerByUserIdHandler)
	// admin.POST(config.MasterDataTrainers, t.authMiddleware.RequireToken("admin"), t.createHandle)         //harusnya bisa dicover pakai csv
	admin.GET(config.MasterDataTrainers, t.authMiddleware.RequireToken("admin"), t.listHandler)           //bisa
	admin.GET(config.MasterDataTrainerByID, t.authMiddleware.RequireToken("admin"), t.trainerByIdHandler) //bisa
	admin.PUT(config.MasterDataTrainerByID, t.authMiddleware.RequireToken("admin", "trainer"), t.updatedTrainerHandler)
	admin.DELETE(config.MasterDataTrainerByID, t.authMiddleware.RequireToken("admin"), t.deleteHandler)
}

func NewTrainerController(trainerUc usecase.TrainerUsecase, rg *gin.RouterGroup, auth middleware.AuthMiddleware) *TrainerController {
	return &TrainerController{
		trainerUc:      trainerUc,
		rg:             rg,
		authMiddleware: auth,
	}
}
