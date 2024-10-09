package auth

import (
	"errors"
	"os"

	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
)

type JWT struct{}

func GetJWTClaims[T any]() (T, error) {
	var claims T
	claimsStr := os.Getenv("JWT_CLAIMS")
	if claimsStr == "" {
		return claims, errors.New("JWT claims not found")
	}
	err := utils.JsonDeserialize([]byte(claimsStr), &claims)
	if err != nil {
		return claims, err
	}

	return claims, nil
}
