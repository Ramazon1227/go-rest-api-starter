package jwt

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Ramazon1227/go-rest-api-starter/models"
)

var (
	SigningKey      []byte
	ExpiryDuration = 24 * time.Hour
)

func GenerateToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.Id
	claims["role"] = user.Role
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ExpiryDuration).Unix()
	claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(SigningKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SigningKey, nil
	})
}

func ExtractClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func GetTokenExpiryTime() time.Time {
	return time.Now().Add(ExpiryDuration)
}

// For token invalidation (blacklisting)
var (
	blacklistedTokens = make(map[string]bool)
	blacklistMu       sync.RWMutex
)

func InvalidateToken(token string) error {
	blacklistMu.Lock()
	defer blacklistMu.Unlock()
	blacklistedTokens[token] = true
	return nil
}

func IsTokenBlacklisted(token string) bool {
	blacklistMu.RLock()
	defer blacklistMu.RUnlock()
	return blacklistedTokens[token]
}
