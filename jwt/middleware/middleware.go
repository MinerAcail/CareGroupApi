package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kobbi/vbciapi/jwt/helpers"
)

// Now, let's define the leaderIDContextKey to use as the context key for LeaderID.
type contextKey string

const LeaderIDContextKey contextKey = "leaderID"
const LeaderTypeContextKey contextKey = "leaderTYPE"

// ValidateToken validates the provided JWT token and returns the claims if the token is valid.
func ValidateTokens(signedToken string) (*helpers.TokenStruct, error) {
	// Parse the JWT token with claims and validation function.
	token, err := jwt.ParseWithClaims(
		signedToken,
		&helpers.TokenStruct{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method used in the token.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return helpers.JwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Verify if the token is valid and not expired.
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract and type-assert the claims.
	claims, ok := token.Claims.(*helpers.TokenStruct)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Verify if the token has expired.
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, errors.New("token has expired")
	}

	// Token is valid, return the claims.
	return claims, nil
}

// AuthenticationMiddleware is the middleware for authentication.
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			next.ServeHTTP(w, r)
			return
		}

		//validate jwt token
		tokenStr := tokenString
		leaderID, leaderType, err := ParseToken(tokenStr)
		if err != nil {
			// Send a JSON-formatted error response.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			response := map[string]string{"error": "Unauthorized token: " + err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		

		// Pass the LeaderID to the context for use in other handlers.
		ctx := context.WithValue(r.Context(), LeaderIDContextKey, leaderID)
		ctx = context.WithValue(ctx, LeaderTypeContextKey, leaderType)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
func ParseToken(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return helpers.JwtSecret, nil
	})
	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		leaderID, ok := claims["leaderID"].(string)
		if !ok {
			return "", "", fmt.Errorf("leaderID not found in the token claims or has an unexpected type")
		}

		leaderType, ok := claims["types"].(string)
		if !ok {
			return "", "", fmt.Errorf("leaderType not found in the token claims or has an unexpected type")
		}

		return leaderID, leaderType, nil
	} else {
		return "", "", fmt.Errorf("invalid token")
	}
}
func ExtractCTXinfo(ctx context.Context) error {
	_, ok := ctx.Value(LeaderIDContextKey).(string)
	if !ok {
		return fmt.Errorf("leaderID not found in request context")
	}
	leaderType, ok := ctx.Value(LeaderTypeContextKey).(string)
	if !ok {
		return fmt.Errorf("leaderType not found in request context")
	}

	leaderType = strings.ToUpper(leaderType) // Convert leaderType to uppercase

	// Check if the leaderType is not one of the allowed values
	if leaderType != "ADMIN" && leaderType != "SUBLEADER" && leaderType != "LEADER" {
		return fmt.Errorf("%s is not allowed", leaderType)
	}

	return nil
}
