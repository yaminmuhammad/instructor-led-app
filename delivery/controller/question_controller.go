package controller

import (
	"fmt"
	"instructor-led-app/config"
	"instructor-led-app/delivery/middleware"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/shared/common"
	"instructor-led-app/usecase"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type QuestionController struct {
	questionUC     usecase.QuestionUseCase
	scheduleUC     usecase.ScheduleUseCase
	trainerUC      usecase.TrainerUsecase
	participantUC  usecase.ParticipantUseCase
	userUC         usecase.UserUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (q *QuestionController) listHandler(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	questions, paging, err := q.questionUC.FindAllQuestion(page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var response []interface{}
	for _, v := range questions {
		response = append(response, v)
	}
	common.SendPagedResponse(ctx, response, paging, "Ok")
}

func (q *QuestionController) UpdatadStatusQuestionByTrainer(ctx *gin.Context) {
	logger := logrus.New()
	var payload dto.QuestionDTO
	user := ctx.MustGet("userID").(string)
	trainer, _ := q.trainerUC.FindTrainerByUserId(user)
	logger.Infoln(trainer.ID)
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	logger.Infoln(payload.ParticipantName)
	userId, _ := q.userUC.FindUserIDByName(payload.ParticipantName)
	logger.Infoln(userId.Id)
	fmt.Println(userId.Id)
	participantId, _ := q.participantUC.GetParticipantByUserId(userId.Id)
	logger.Infoln(participantId.ID)
	fmt.Println(participantId.ID)
	UpdatedQuestion, err := q.questionUC.UpdatadStatusQuestionByTrainer(trainer.ID, participantId.ID, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed update"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": UpdatedQuestion})
}

func (q *QuestionController) CreateQuestionByTrainer(ctx *gin.Context) {
	var payload dto.QuestionDTO
	user := ctx.MustGet("userID").(string)
	trainer, _ := q.trainerUC.FindTrainerByUserId(user)
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, _ := q.userUC.FindUserIDByName(payload.ParticipantName)
	participantId, _ := q.participantUC.GetParticipantByUserId(userId.Id)
	InsertQuestion, err := q.questionUC.CreateQuestionByTrainer(trainer.ID, participantId.ID, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed update"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": InsertQuestion})
}

func (q *QuestionController) createHandler(ctx *gin.Context) {
	var payload entity.Question
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	question, err := q.questionUC.CreateNewQuestion(payload)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	common.SendCreateResponse(ctx, question, "Created")
}

func (q *QuestionController) deleteHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := q.questionUC.DeleteQuestion(id)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	common.SendErrorResponse(ctx, http.StatusNoContent, "Deleted Successfully")
}

func (q *QuestionController) update(c *gin.Context) {
	id := c.Param("id")
	var payload entity.Question
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	question, err := q.questionUC.UpdateQuestion(id, payload)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.SendSingleResponse(c, question, "Updated Successfully")
}

func (q *QuestionController) getById(c *gin.Context) {
	id := c.Param("id")
	question, err := q.questionUC.FindById(id)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.SendSingleResponse(c, question, "Ok")
}

func (q *QuestionController) GetQuestionByTrainerId(c *gin.Context) {
	// id := c.Param("id")
	userId := c.MustGet("user").(string)

	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	question, paging, err := q.questionUC.FindQuestionByTrainerId(userId, page, size)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var response []interface{}
	for _, v := range question {
		response = append(response, v)
	}
	common.SendPagedResponse(c, response, paging, "Ok")
}

func (q *QuestionController) NewQuestionByPartcipant(ctx *gin.Context) {
	var payload dto.QuestionDto

	// Menggunakan ID peserta dari konteks
	userId := ctx.MustGet("userID").(string)

	participantID, err := q.participantUC.GetParticipantByUserId(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get participantID"})
		return
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Buat pertanyaan tanpa memasukkan participantId dari payload
	newQuestion, err := q.questionUC.CreateQuestionByParticipant(participantID.ID, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newQuestion})
}
func (q *QuestionController) Route() {
	q.rg.GET(config.QuestionGetList, q.authMiddleware.RequireToken("participant"), q.listHandler)
	q.rg.GET(config.QuestionGetById, q.authMiddleware.RequireToken("participant", "trainer"), q.getById)
	q.rg.POST(config.QuestionPost, q.authMiddleware.RequireToken("participant"), q.createHandler)
	q.rg.PUT(config.QuestionGetById, q.authMiddleware.RequireToken("trainer"), q.update)
	q.rg.DELETE(config.QuestionDelete, q.authMiddleware.RequireToken("admin"), q.deleteHandler)
	q.rg.GET(config.QuestionTrainer, q.authMiddleware.RequireToken("trainer"), q.GetQuestionByTrainerId)
	q.rg.POST(config.QuestionTrainer, q.authMiddleware.RequireToken("trainer"), q.CreateQuestionByTrainer)
	q.rg.PUT(config.UpdateQuestionByTrainer, q.authMiddleware.RequireToken("trainer"), q.UpdatadStatusQuestionByTrainer)
	q.rg.POST(config.ParticipantNewQuetion, q.authMiddleware.RequireToken("participant"), q.NewQuestionByPartcipant)
}
func NewQuestionController(questionUC usecase.QuestionUseCase, scheduleUC usecase.ScheduleUseCase, trainerUC usecase.TrainerUsecase, participantUC usecase.ParticipantUseCase, userUC usecase.UserUsecase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *QuestionController {
	return &QuestionController{
		questionUC:     questionUC,
		trainerUC:      trainerUC,
		scheduleUC:     scheduleUC,
		participantUC:  participantUC,
		userUC:         userUC,
		rg:             rg,
		authMiddleware: authMiddleware,
	}
}
