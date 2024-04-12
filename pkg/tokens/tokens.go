package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrTokenExpired = errors.New("token is expired")
	ErrTokenInvalid = errors.New("token is invalid")
)

type TokenManagerInterface interface {
	CreateAccessToken(userId uuid.UUID) (string, error)
	CreateEmailToken(email string) (string, error)
	ParseAccessToken(tokenString string) (uuid.UUID, error)
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

type ClaimsAccessToken struct {
	jwt.StandardClaims
	Id uuid.UUID `json:"id"`
}

type ClaimsEmailToken struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

type ClaimsMachineToken struct {
	jwt.StandardClaims
	Id        uuid.UUID `json:"id"`
	IsMachine struct{}  `json:"isMachine"`
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

func (tm *TokenManager) CreateAccessToken(userId uuid.UUID) (string, error) {
	return tm.createJWTToken(&ClaimsAccessToken{
		tm.createStandartClaims(tm.accessTTL),
		userId,
	})
}

// func (tm *TokenManager) CreateMachineToken(machineId uuid.UUID) (string, error) {
// 	return tm.createJWTToken(ClaimsMachineToken{
// 		tm.createStandartClaims(time.Hour * 24 * 365 * 100),
// 		machineId,
// 		struct{}{},
// 	})
// }

func (tm *TokenManager) CreateEmailToken(email string) (string, error) {
	return tm.createJWTToken(&ClaimsEmailToken{
		tm.createStandartClaims(tm.emailTTL),
		email,
	})
}

func (tm *TokenManager) ParseAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsAccessToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Println("error invalid")
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
		logrus.Errorf("[tokens]: error parsing access token: %s", err)
		if errors.Is(ErrTokenExpired, err) {
			return uuid.UUID{}, ErrTokenExpired
		}
		return uuid.UUID{}, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*ClaimsAccessToken)
	if !ok {
		return uuid.UUID{}, ErrTokenInvalid
	}

	return claims.Id, nil
}

func (tm *TokenManager) ParseEmailToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsEmailToken{}, func(t *jwt.Token) (interface{}, error) {
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
		logrus.Errorf("[tokens]: error parsing email token: %s", err)
		if errors.Is(ErrTokenExpired, err) {
			return "", ErrTokenExpired
		}
		return "", ErrTokenInvalid
	}
	claims, ok := token.Claims.(*ClaimsEmailToken)
	if !ok {
		return "", ErrTokenInvalid
	}

	return claims.Email, nil
}

// func (tm *TokenManager) ParseMachineToken(tokenString string) (uuid.UUID, error) {
// 	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
// 		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, ErrTokenInvalid
// 		}
// 		claims, ok := t.Claims.(*ClaimsMachineToken)
// 		if !ok {
// 			return nil, ErrTokenInvalid
// 		}
// 		if claims.ExpiresAt <= time.Now().Unix() {
// 			return nil, ErrTokenExpired
// 		}
// 		return []byte(tm.secretKey), nil
// 	})

// 	if err != nil {
// 		return uuid.UUID{}, err
// 	}

// 	claims, ok := token.Claims.(*ClaimsMachineToken)
// 	if !ok {
// 		return uuid.UUID{}, ErrTokenInvalid
// 	}

// 	return claims.Id, nil
// }
