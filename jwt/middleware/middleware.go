package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kobbi/vbciapi/jwt/helpers"
)

// Now, let's define the IDContextKey to use as the context key for LeaderID.
type contextKey string

const IDContextKey contextKey = "leaderID"
const SubIDContextKey contextKey = "subsID"
const IDoptionalContextKey contextKey = "optional1"
const LeaderTypeContextKey contextKey = "leaderTYPE"
const cookieAccessKeyCtx contextKey = "cookieAccessKeyCtx"
const CookieName contextKey = "SetToken"

// ValidateToken validates the provided JWT token and returns the claims if the token is valid.

// CookieAccess struct for managing cookies
type CookieAccess struct {
	Writer     http.ResponseWriter
	UserId     string
	IsLoggedIn bool
}

// SetToken method to write a cookie
func (ca *CookieAccess) SetToken(token string) {
	http.SetCookie(ca.Writer, &http.Cookie{
		Name:     string(CookieName),
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour), // Adjust expiration as needed

	})
}
func (ca *CookieAccess) Logout() error {
	cookie := &http.Cookie{
		Name:     string(CookieName),
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0), // Expire immediately
	}

	http.SetCookie(ca.Writer, cookie)

	return nil
}

// Function to set CookieAccess object in the request context
func setValInCtx(ctx context.Context, val interface{}) context.Context {
	return context.WithValue(ctx, cookieAccessKeyCtx, val)
}

// Function to retrieve CookieAccess object from context
func GetCookieAccess(ctx context.Context) *CookieAccess {
	return ctx.Value(cookieAccessKeyCtx).(*CookieAccess)
}

// Middleware function for Chi
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieA := CookieAccess{
			Writer: w,
		}

		ctx := r.Context()

		// Extract and parse the token here
		tokenCookie, err := r.Cookie(string(CookieName))
		if err == nil {
			token := tokenCookie.Value
			leaderID, leaderType, optionalClaims, err := ParseToken(token)
			if err == nil {
				ctx = context.WithValue(ctx, IDContextKey, leaderID)
				ctx = context.WithValue(ctx, LeaderTypeContextKey, leaderType)

				// Set optional claims in the context using custom keys
				for key, value := range optionalClaims {
					ctx = context.WithValue(ctx, contextKey(key), value)
				}
			}
		}

		ctx = setValInCtx(ctx, &cookieA)

		// Extract user ID from cookie or perform any necessary checks here

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ParseToken(tokenStr string) (string, string, map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return helpers.JwtSecret, nil
	})
	if err != nil {
		return "", "", nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		leaderID, ok := claims["leaderID"].(string)
		if !ok {
			return "", "", nil, fmt.Errorf("leaderID not found in the token claims or has an unexpected type")
		}

		leaderType, ok := claims["types"].(string)
		if !ok {
			return "", "", nil, fmt.Errorf("types not found in the token claims or has an unexpected type")
		}

		// Extract optional claims, if they exist
		optionalClaims := make(map[string]interface{})
		for key, value := range claims {
			if key != "leaderID" && key != "types" {
				optionalClaims[key] = value
			}
		}

		return leaderID, leaderType, optionalClaims, nil
	} else {
		return "", "", nil, fmt.Errorf("invalid token")
	}
}

func ExtractCTXinfo(ctx context.Context) error {
	_, ok := ctx.Value(IDContextKey).(string)
	if !ok {
		return fmt.Errorf("ID not found in request context")
	}
	leaderType, ok := ctx.Value(LeaderTypeContextKey).(string)
	if !ok {
		return fmt.Errorf("leaderType not found in request context")
	}

	leaderType = strings.ToUpper(leaderType) // Convert leaderType to uppercase

	// Check if the leaderType is not one of the allowed values
	if leaderType != "ADMIN"&& leaderType != "SUBCHURCH" && leaderType != "SUBLEADER" && leaderType != "LEADER" && leaderType != "CALLAGENT" {
		return fmt.Errorf("%s is not allowed", leaderType)
	}
	
	
	return nil
}
func ExtractCTXinfo4AdminOnly(ctx context.Context) error {
	_, ok := ctx.Value(IDContextKey).(string)
	if !ok {
		return fmt.Errorf("leader not found in request context")
	}
	leaderType, ok := ctx.Value(LeaderTypeContextKey).(string)
	if !ok {
		return fmt.Errorf("leader Type not found in request context")
	}

	leaderType = strings.ToUpper(leaderType) // Convert leaderType to uppercase

	// Check if the leaderType is not one of the allowed values
	if leaderType != "ADMIN" && leaderType != "SUBCHURCH" {
		return fmt.Errorf("%s is not allowed", leaderType)
	}

	return nil
}

func ExtractCTXinfo4CallCenter(ctx context.Context) error {
	_, ok := ctx.Value(IDContextKey).(string)
	if !ok {
		return fmt.Errorf("leader not found in request context")
	}
	leaderType, ok := ctx.Value(LeaderTypeContextKey).(string)
	if !ok {
		return fmt.Errorf("leader Type not found in request context")
	}

	leaderType = strings.ToUpper(leaderType) // Convert leaderType to uppercase

	// Check if the leaderType is not one of the allowed values
	if leaderType != "CALLCENTER"{
		return fmt.Errorf("%s is not allowed", leaderType)
	}

	return nil
}

// AuthenticationMiddleware is the middleware for authentication.
/* func AuthenticationMiddleware(next http.Handler) http.Handler {
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
		ctx := context.WithValue(r.Context(), IDContextKey, leaderID)
		ctx = context.WithValue(ctx, LeaderTypeContextKey, leaderType)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
} */

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
	if claims.ExpiresAt> time.Now().UTC().Unix() {
		return nil, errors.New("token has expired")
	}

	// Token is valid, return the claims.
	return claims, nil
}
