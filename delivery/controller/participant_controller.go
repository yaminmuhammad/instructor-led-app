package controller

import (
	"instructor-led-app/config"
	"instructor-led-app/delivery/middleware"
	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/common"
	"instructor-led-app/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type participantController struct {
	participantUseCase usecase.ParticipantUseCase
	userUc             usecase.UserUsecase
	rg                 *gin.RouterGroup
	authMiddleware     middleware.AuthMiddleware
}

func (c *participantController) insertHandler(ctx *gin.Context) {
	var participantDto dto.ParticipantDTO
	if err := ctx.ShouldBindJSON(&participantDto); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	participantDto, err := c.participantUseCase.CreateNewParticipant(participantDto)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	common.SendCreateResponse(ctx, participantDto, "Create participant successfully")
}

func (c *participantController) listHandler(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))

	participants, paging, err := c.participantUseCase.GetAllParticipants(page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var response []interface{}
	for _, value := range participants {
		response = append(response, value)
	}

	common.SendPagedResponse(ctx, response, paging, "Get all Participants successfully")
}

func (c *participantController) getHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	participant, err := c.participantUseCase.GetParticipantByID(id)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	common.SendSingleResponse(ctx, participant, "Get Participant successfully")
}

//	 Warn: tidak digunakan
// func (c *participantController) getScheduleByParticipantIDHandler(ctx *gin.Context) {
// 	user := ctx.MustGet("userID").(string)

// 	participant, err := c.participantUseCase.FindScheduleWithParticipantId(user)
// 	fmt.Println(participant)
// 	if err != nil {
// 		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	common.SendSingleResponse(ctx, participant, "Get Participant successfully")

// }

// // Deprecated
// func (c *participantController) updateHandler(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	var participantDto dto.ParticipantDTO
// 	if err := ctx.ShouldBindJSON(&participantDto); err != nil {
// 		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	participantDto.ID = id

// 	participant, err := c.participantUseCase.UpdateParticipantByID(participantDto)
// 	if err != nil {
// 		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	common.SendSingleResponse(ctx, participant, "Update Participant successfully")
// }

func (c *participantController) deleteHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.participantUseCase.DeleteParticipantByID(id); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	common.SendDeleteResponse(ctx, "Delete Participant successfully")
}

func (c *participantController) UpdateParticipantRoleByAdmin(ctx *gin.Context) {
	var payload dto.ParticipantRoleDTO

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		return
	}
	userId, err := c.userUc.FindUserIDByName(payload.Name)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "user not found")
		return
	}

	participantId, err := c.participantUseCase.GetParticipantByUserId(userId.Id)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "participant not found")
		return
	}

	err = c.participantUseCase.UpdateParticipantByRole(payload.Role, participantId.ID)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to update Participant: "+err.Error())
		return
	}

	common.SendSingleResponse(ctx, err, "Ok")

}

func (c *participantController) UpdateParticipantByParticipantId(ctx *gin.Context) {
	user := ctx.MustGet("userID").(string)
	participantId, err := c.participantUseCase.GetParticipantByUserId(user)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "invalid ID")
		return
	}
	var participant dto.ParticipantDTO

	if err := ctx.ShouldBindJSON(&participant); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		return
	}

	participant.ID = participantId.ID
	updateParticipant, err := c.participantUseCase.UpdateParticipantByID(participant)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to update trainer: "+err.Error())
		return
	}
	common.SendSingleResponse(ctx, updateParticipant, "Ok")
}

func (c *participantController) Route() {
	admin := c.rg.Group(config.AdminGroup)
	admin.POST(config.MasterDataParticipants, c.authMiddleware.RequireToken("admin"), c.insertHandler)           //ini harusnya ditanganin sama csv
	admin.GET(config.MasterDataParticipants, c.authMiddleware.RequireToken("trainer", "admin"), c.listHandler)   //bisa
	admin.GET(config.MasterDataParticipantByID, c.authMiddleware.RequireToken("trainer", "admin"), c.getHandler) //bisa
	admin.PUT(config.MasterDataParticipants, c.authMiddleware.RequireToken("participant", "admin"), c.UpdateParticipantByParticipantId)
	admin.PUT(config.MasterDataParticipantsRole, c.authMiddleware.RequireToken("admin"), c.UpdateParticipantRoleByAdmin) //bisa
	admin.DELETE(config.MasterDataParticipantByID, c.authMiddleware.RequireToken("admin"), c.deleteHandler)
}

func NewParticipantController(participantUseCase usecase.ParticipantUseCase, userUC usecase.UserUsecase, rg *gin.RouterGroup, auth middleware.AuthMiddleware) *participantController {
	return &participantController{participantUseCase, userUC, rg, auth}
}
