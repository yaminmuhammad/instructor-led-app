package usecase

import (
	"errors"
	"fmt"
	"instructor-led-app/entity"
	"instructor-led-app/entity/dto"
	"instructor-led-app/repository"
	"instructor-led-app/shared/model"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	FindById(id string) (entity.User, error)
	FindAllUser(page, size int) ([]entity.User, model.Paging, error)
	CreatedUser(data entity.User) (entity.User, error)
	CreatedUserByCsv(filePath string) ([]entity.User, error)
	UpdatedUser(id string, data entity.User) (entity.User, error)
	DeleteUser(id string) (entity.User, error)
	AuthUser(email string, hashPassword string) (entity.User, error)
	FindUserIDByName(name string) (dto.UserId, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

// FindUserIDByName implements UserUsecase.
func (u *userUsecase) FindUserIDByName(name string) (dto.UserId, error) {
	return u.repo.GetUserIDByName(name)
}

// AuthUser implements UserUsecase.
func (t *userUsecase) AuthUser(email string, hashPassword string) (entity.User, error) {
	user, err := t.repo.GetUser(email)
	if err != nil {
		return entity.User{}, err
	}

	if hashPassword == "" {
		log.Println("userUsecase.AuthUser: Empty password provided")
		return entity.User{}, errors.New("password required")
	}

	// Verify the entered password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Hashpassword), []byte(hashPassword))
	if err != nil {
		log.Println("userUsecase.AuthUser: Password verification failed")
		return entity.User{}, errors.New("password verification failed")

	}

	return user, nil

}

// FindAllUser implements UserUsecase.
func (t *userUsecase) FindAllUser(page int, size int) ([]entity.User, model.Paging, error) {
	return t.repo.List(page, size)

}

// CreatedUserByCsv implements UserUsecase.
func (t *userUsecase) CreatedUserByCsv(filePath string) ([]entity.User, error) {
	var users []entity.User

	// Baca data dari file Csv menggunakan repository
	users, err := t.repo.CreateByCsv(filePath)
	if err != nil {
		log.Println("userUseCase.CreatedUserByCsv: Error creating users from CSV:", err.Error())
		return nil, err
	}
	return users, nil
}

// CreatedUser implements UserUsecase.
func (t *userUsecase) CreatedUser(data entity.User) (entity.User, error) {

	if data.Name == "" || data.Email == "" || data.Username == "" || data.Address == "" || data.Hashpassword == "" || data.Role == "" {
		return entity.User{}, fmt.Errorf("oppps, required fields")
	}
	data.UpdatedAt = time.Now()
	user, err := t.repo.Created(data)
	if err != nil {
		return entity.User{}, fmt.Errorf("oppps, failed to save data user :%v", err.Error())
	}
	return user, nil

}

// DeleteCustomer implements UserUsecase.
func (t *userUsecase) DeleteUser(id string) (entity.User, error) {
	return t.repo.Delete(id)

}

// FindById implements userUsecase.
func (t *userUsecase) FindById(id string) (entity.User, error) {
	return t.repo.Get(id)
}

// UpdatedCustomer implements UserUsecase.
func (t *userUsecase) UpdatedUser(id string, data entity.User) (entity.User, error) {
	data.UpdatedAt = time.Now()
	data, err := t.repo.Updated(id, data)
	if err != nil {
		return entity.User{}, fmt.Errorf("oppps, failed to update data user :%v", err.Error())
	}
	return data, nil

}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}
