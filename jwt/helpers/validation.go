package helpers

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
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

func GenerateToken(types string, leaderID string, optionalClaims ...jwt.MapClaims) (string, error) {
	// Create a new JWT token with the signing method HS256 (HMAC SHA-256)
	token := jwt.New(jwt.SigningMethodHS256)

	// Convert the token's claims to a map
	claims := token.Claims.(jwt.MapClaims)

	// Set the "types" and "leaderID" claims in the token
	claims["types"] = types
	claims["leaderID"] = leaderID

	// Set the "exp" (expiration) claim to the current time + 24 hours
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Set optional claims if provided
	for _, optionalClaim := range optionalClaims {
		for key, value := range optionalClaim {
			claims[key] = value
		}
	}

	// Sign the token using a secret key (SecretKey)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}

	// Return the signed token as a string
	return tokenString, nil
}

func GenerateRandomPassword(length int) string {
   // Define a character set for the password
   charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
   password := make([]byte, length)

   // Generate the password
   for i := range password {
       password[i] = charset[rand.Intn(len(charset))]
   }

   return string(password)
}

func GenerateRandomEmail(domain string,username string) string {
	// username := GenerateRandomPassword(6) // Generate a random username
	return fmt.Sprintf("%s@%s", username, domain)
}

// func GenerateTokenO(leaderID string, optionalClaims ...jwt.MapClaims) (string, error) {
// 	// Create a new JWT token with the signing method HS256 (HMAC SHA-256)
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	// Convert the token's claims to a map
// 	claims := token.Claims.(jwt.MapClaims)

// 	// Set the "leaderID" claim in the token
// 	claims["leaderID"] = leaderID

// 	// Set any optional claims provided
// 	if len(optionalClaims) > 0 {
// 		for key, value := range optionalClaims[0] {
// 			claims[key] = value
// 		}
// 	}

// 	// Set the "exp" (expiration) claim to a reasonable future time (e.g., 24 hours)
// 	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

// 	// Sign the token using a secret key (JwtSecret)
// 	tokenString, err := token.SignedString([]byte(JwtSecret))
// 	if err != nil {
// 		return "", err
// 	}

// 	// Return the signed token as a string
// 	return tokenString, nil
// }
