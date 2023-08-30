package helpers

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)
var JwtSecret = []byte("happytime") // Replace this with your own secret key

type TokenStruct struct {
	ID    string `json:"id"`
	Types string `json:"types"`
	jwt.StandardClaims
}


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VerifyPassword(hashedPassword string, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		return errors.New("email or password is incorrect")
	}
	return nil
}

func GenerateTokens(types string, id string) (string, error) {
	// Access Token

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&TokenStruct{
			ID:    id,
			Types: types,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
				Id:        id,
			},
		})

	signedToken, err := accessToken.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}

	// Refresh Token
	/* refreshTokenClaims := &TokenStruct{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(JwtSecret))
	if err != nil {
		return "", "", err
	} */

	return signedToken, nil
}
func GenerateToken(types string, leaderID string) (string, error) {
	// Create a new JWT token with the signing method HS256 (HMAC SHA-256)
	token := jwt.New(jwt.SigningMethodHS256)

	// Convert the token's claims to a map
	claims := token.Claims.(jwt.MapClaims)

	// Set the "types" and "leaderID" claims in the token
	claims["types"] = types
	claims["leaderID"] = leaderID

	// Set the "exp" (expiration) claim to the current time + 24 hours
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Sign the token using a secret key (SecretKey)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}

	// Return the signed token as a string
	return tokenString, nil
}
