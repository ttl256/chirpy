package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrMissingAuthHeader = errors.New("missing auth header")
var ErrMalformedAuthHeader = errors.New("malformed auth header")

func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(id int, secret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(id),
		Audience:  []string{},
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		NotBefore: &jwt.NumericDate{
			Time: time.Time{},
		},
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ID:       "",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(token, secret string) (string, error) {
	t, err := jwt.ParseWithClaims(
		token,
		&jwt.RegisteredClaims{}, //nolint: exhaustruct // let me be
		func(_ *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	if err != nil {
		return "", err
	}

	subject, err := t.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return subject, nil
}

func GetAuth(header http.Header, key string) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}
	hFileds := strings.Fields(authHeader)
	if len(hFileds) < 2 || hFileds[0] != key {
		return "", ErrMalformedAuthHeader
	}
	return hFileds[1], nil
}

func RandString64(n int) (string, error) {
	xs := make([]byte, n)
	_, err := rand.Read(xs)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(xs), nil
}
