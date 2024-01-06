package controller

import (
	"instructor-led-app/config"
	"instructor-led-app/delivery/middleware"
	"instructor-led-app/entity"
	"instructor-led-app/shared/common"
	"instructor-led-app/usecase"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUC         usecase.UserUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (t *UserController) ListHandler(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))

	user, paging, err := t.userUC.FindAllUser(page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	var response []interface{}
	for _, v := range user {
		response = append(response, v)
	}
	common.SendPagedResponse(ctx, response, paging, "Ok")

}

func (t *UserController) createdByCsv(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "Error getting file: "+err.Error())
		return
	}

	// Simpan file CSV yang diunggah ke server
	filePath := "uploads/csv/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "Error saving file: "+err.Error())
		return
	}

	// Panggil use case untuk menangani pembacaan dan penyimpanan data dari file CSV
	users, err := t.userUC.CreatedUserByCsv(filePath)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "Error creating users from CSV: "+err.Error())
		return
	}

	// Hapus file CSV yang telah diunggah
	if err := os.Remove(filePath); err != nil {
		log.Println("Error removing uploaded file:", err.Error())
	}

	common.SendSingleResponse(c, users, "Users created successfully from CSV")
}

func (t *UserController) getById(c *gin.Context) {
	id := c.Param("id")
	user, err := t.userUC.FindById(id)
	if err != nil {
		common.SendErrorResponse(c, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}
	common.SendSingleResponse(c, user, "ok")
}

func (t *UserController) create(c *gin.Context) {
	var data entity.User
	if err := c.ShouldBindJSON(&data); err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := t.userUC.CreatedUser(data)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.SendSingleResponse(c, user, "Created")
}

func (t *UserController) update(c *gin.Context) {
	id := c.Param("id")
	// Pengecekan ID valid
	if id == "" {
		common.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	// Pencarian customer berdasarkan ID
	_, err := t.userUC.FindById(id)
	if err != nil {
		common.SendErrorResponse(c, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}
	var data entity.User
	// Pengecekan dan binding data JSON
	if err := c.ShouldBindJSON(&data); err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, "Invalid JSON data: "+err.Error())
		return
	}
	// Pembaruan user
	user, err := t.userUC.UpdatedUser(id, data)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update customer: "+err.Error())
		return
	}
	common.SendSingleResponse(c, user, "Updated successfully")
}

func (t *UserController) delete(c *gin.Context) {
	id := c.Param("id")
	// Pengecekan ID valid
	if id == "" {
		common.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	// Pencarian user berdasarkan ID
	_, err := t.userUC.FindById(id)
	if err != nil {
		common.SendErrorResponse(c, http.StatusNotFound, "Customer with ID "+id+" not found")
		return
	}
	_, err = t.userUC.DeleteUser(id)
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete customer: "+err.Error())
		return
	}
	common.SendSingleResponse(c, "Ok", "Delete successfully")

}

func (t *UserController) Route() {
	admin := t.rg.Group(config.AdminGroup)
	admin.POST(config.MasterDataUsersCsv, t.authMiddleware.RequireToken("admin"), t.createdByCsv) //bisa
	admin.POST(config.MasterDataUsers, t.create)                                                  //bisa
	// admin.POST(config.MasterDataUsers, t.authMiddleware.RequireToken("admin"), t.create)          //bisa
	admin.GET(config.MasterDataUsers, t.authMiddleware.RequireToken("admin"), t.ListHandler)  //bisa
	admin.GET(config.MasterDataUserByID, t.authMiddleware.RequireToken("admin"), t.getById)   //bisa
	admin.PUT(config.MasterDataUserByID, t.authMiddleware.RequireToken("admin"), t.update)    //bisa
	admin.DELETE(config.MasterDataUserByID, t.authMiddleware.RequireToken("admin"), t.delete) //bisa
}

func NewUserController(userUC usecase.UserUsecase, rg *gin.RouterGroup, auth middleware.AuthMiddleware) *UserController {
	return &UserController{
		userUC:         userUC,
		rg:             rg,
		authMiddleware: auth,
	}
}
