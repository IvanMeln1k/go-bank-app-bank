package tokens

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrTokenExpired = errors.New("token is expired")
	ErrTokenInvalid = errors.New("token is invalid")
)

type TokenManagerInterface interface {
	CreateRefreshToken() (string, error)
	CreateAccessToken(userId int) (string, error)
	CreateEmailToken(email string) (string, error)
	ParseAccessToken(tokenString string) (int, error)
	ParseEmailToken(tokenString string) (string, error)
}

type TokenManager struct {
	secretKey string
	accessTTL time.Duration
	emailTTL  time.Duration
}

type Config struct {
	SecretKey string
	AccessTTL time.Duration
	EmailTTL  time.Duration
}

func NewTokenManager(cfg Config) *TokenManager {
	return &TokenManager{
		secretKey: cfg.SecretKey,
		accessTTL: cfg.AccessTTL,
		emailTTL:  cfg.EmailTTL,
	}
}

func (tm *TokenManager) CreateRefreshToken() (string, error) {
	b := make([]byte, 32)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

type ClaimsAccessToken struct {
	jwt.StandardClaims
	Id int `json:"id"`
}

type ClaimsEmailToken struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func (tm *TokenManager) createStandartClaims(ttl time.Duration) jwt.StandardClaims {
	return jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
}

func (tm *TokenManager) createJWTToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.secretKey))
}

func (tm *TokenManager) CreateAccessToken(userId int) (string, error) {
	return tm.createJWTToken(ClaimsAccessToken{
		tm.createStandartClaims(tm.accessTTL),
		userId,
	})
}

func (tm *TokenManager) CreateEmailToken(email string) (string, error) {
	return tm.createJWTToken(ClaimsEmailToken{
		tm.createStandartClaims(tm.emailTTL),
		email,
	})
}

func (tm *TokenManager) ParseAccessToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		claims, ok := t.Claims.(*ClaimsAccessToken)
		if !ok {
			return nil, ErrTokenInvalid
		}
		if claims.ExpiresAt <= time.Now().Unix() {
			return nil, ErrTokenExpired
		}
		return []byte(tm.secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*ClaimsAccessToken)
	if !ok {
		return 0, ErrTokenInvalid
	}

	return claims.Id, nil
}

func (tm *TokenManager) ParseEmailToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		claims, ok := t.Claims.(*ClaimsEmailToken)
		if !ok {
			return nil, ErrTokenInvalid
		}
		if claims.ExpiresAt <= time.Now().Unix() {
			return nil, ErrTokenExpired
		}
		return []byte(tm.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*ClaimsEmailToken)
	if !ok {
		return "", ErrTokenInvalid
	}

	return claims.Email, nil
}
