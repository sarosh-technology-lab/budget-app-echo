package services

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/common"
	"budget-backend/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (userService UserService) RegisterUser(userRequest *requests.RegisterUserRequest) (*models.User, error) {
	// hash password
	hashedPassword, err := common.HashPassword(userRequest.Password)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("User")
	}
	user := models.User{
		FirstName: userRequest.FirstName,
		LastName: userRequest.LastName,
		Password: hashedPassword,
		Email: userRequest.Email,
	}

	result := userService.db.Create(&user)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, errors.New("user registration failed")
	}
	return &user, nil
}

func (userService UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := userService.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}