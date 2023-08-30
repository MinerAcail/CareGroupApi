package helpers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	// "github.com/99designs/gqlgen/graphql"
	"github.com/dgrijalva/jwt-go"
)

// getTokenFromRequest extracts the JWT token from the Authorization header in the HTTP request.
func getTokenFromRequest(ctx context.Context) string {
	req, _ := ctx.Value("request").(*http.Request)
	if req == nil {
		return ""
	}

	authorizationHeader := req.Header.Get("Authorization")
	if authorizationHeader == "" {
		return ""
	}

	// The token should be in the format: "Bearer <token>"
	splitHeader := strings.Split(authorizationHeader, " ")
	if len(splitHeader) != 2 || strings.ToLower(splitHeader[0]) != "bearer" {
		return ""
	}

	return splitHeader[1]
}

// ExtractLeaderIDFromToken takes a JWT token, verifies it, and extracts the LeaderID from the token's claims.
func ExtractLeaderIDFromToken(tokenString string) (int, error) {
	

	// Parse and validate the token.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token signing method is correct.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return JwtSecret, nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	// Check if the token is valid and not expired.
	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Extract the LeaderID from the token's claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("failed to extract token claims")
	}

	leaderID, ok := claims["leaderID"].(float64)
	if !ok {
		return 0, fmt.Errorf("leaderID not found in token claims")
	}

	return int(leaderID), nil
}


