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
		return nil, errors.New("password hashing failed")
	}
	user := models.User{
		RoleId: userRequest.RoleId,
		FirstName: userRequest.FirstName,
		LastName: userRequest.LastName,
		Password: hashedPassword,
		Phone: userRequest.Phone,
		Email: userRequest.Email,
		Gender: &userRequest.Gender,
		Address: userRequest.Address,
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

func (userService UserService) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	result := userService.db.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (userService UserService) ChangeUserPassword(user models.User, newPassword string) error {
	// hash password
	hashedPassword, err := common.HashPassword(newPassword)
	if err != nil {
		fmt.Println(err)
		return errors.New("password change failed")
	}

	result := userService.db.Model(user).Update("password", hashedPassword); 
	if result.Error != nil {
		return errors.New("password change failed")
	}

	return nil
}

func (userService UserService) UpdateUser(user models.User, updateMap map[string]any) error {
	result := userService.db.Model(user).Updates(updateMap); 
	if result.Error != nil {
		return errors.New("update user failed")
	}

	return nil
}