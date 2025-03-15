package services

import (
	"budget-backend/internal/models"
	"errors"
	"math/rand"
	"strconv"
	"time"
	"gorm.io/gorm"
)

type AppTokenService struct {
	db *gorm.DB
}

func NewAppTokenService(db *gorm.DB) *AppTokenService {
	return &AppTokenService{db: db}
}

func (appTokenService *AppTokenService) getToken() int {
	rand.Seed((time.Now().UnixNano()))
	min := 10000
	max := 99999
	return rand.Intn(max-min+1) + min
}

func (appTokenService *AppTokenService) GenerateResetPasswordToken(user models.User) (*models.AppToken, error) {
	tokenCreated := models.AppToken{
		TargetId:  user.ID,
		Type:      "reset_password",
		Token:     strconv.Itoa(appTokenService.getToken()),
		Used:      false,
		ExpiresAt: time.Now().Add((time.Hour * 24)),
	}
	result := appTokenService.db.Create(&tokenCreated)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tokenCreated, nil
}

func (appTokenService *AppTokenService) ValidatePasswordToken(user models.User, token string) (*models.AppToken, error) {

	var retreivedToken = models.AppToken{}
	result := appTokenService.db.Where(&models.AppToken{
		TargetId: user.ID,
		Type:     "reset_password",
		Token:    token,
	}).First((&retreivedToken))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid passowrd reset token")
		}
		return nil, result.Error
	}

	if retreivedToken.Used {
		return nil, errors.New("reset password token already used")
	}

	if retreivedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("reset password token is expired, please re-initiate forgot password")
	}

	if result.Error != nil {
		return nil, result.Error
	}
	return &retreivedToken, nil
}

func (appTokenService *AppTokenService) InvalidatedToken(user_id int, appToken models.AppToken) {
	appTokenService.db.Model(&models.AppToken{}).Where("trget_id = ? AND token = ?", user_id, appToken.Token).Update("used", true)
}
