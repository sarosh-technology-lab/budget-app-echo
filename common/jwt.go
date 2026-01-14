package common

import (
	"budget-backend/internal/models"
	"errors"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CustomJWTClaims struct{
	ID uint `json:"id"`
	RoleID uint `json:"role_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User, db *gorm.DB) (*string, *string, error) {
    // Access token (short-lived, e.g., 15 minutes)
    userClaims := CustomJWTClaims{
        ID:     uint(user.ID),
        RoleID: user.RoleId,
        RegisteredClaims: jwt.RegisteredClaims{
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 45)), // Short expiration 45 minutes
        },
    }
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
    signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        echo.New().Logger.Error(err)
        return nil, nil, err
    }

    // Refresh token (long-lived, e.g., 7 days)
    refreshClaims := CustomJWTClaims{
        ID:     uint(user.ID),
        RoleID: user.RoleId,
        RegisteredClaims: jwt.RegisteredClaims{
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
        },
    }
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        echo.New().Logger.Error(err)
        return nil, nil, err
    }

    // Store refresh token in the database
    refreshTokenModel := models.RefreshToken{
        UserID:    uint(user.ID),
        Token:     signedRefreshToken,
        ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
    }
    if err := db.Create(&refreshTokenModel).Error; err != nil {
        echo.New().Logger.Error(err)
        return nil, nil, err
    }

    return &signedAccessToken, &signedRefreshToken, nil
}

func ParseJWTSignedAccessToken(signedAccessToken string) (*CustomJWTClaims, error) {
	parsedJWTAccessToken, err := jwt.ParseWithClaims(signedAccessToken, &CustomJWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		echo.New().Logger.Error(err)
		return nil, err
	} else if claims, ok := parsedJWTAccessToken.Claims.(*CustomJWTClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("unknown claims type, cannot proceed")
	}
}

func IsClaimExpired(claims *CustomJWTClaims) bool{
	currentTime := jwt.NewNumericDate(time.Now())
	return claims.ExpiresAt.Time.Before(currentTime.Time)
}

