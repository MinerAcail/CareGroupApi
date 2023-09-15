package middleware

import "net/http"

func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true, // Ensure the cookie is not accessible via JavaScript
		Secure:   true, // Use "Secure" only in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	})
}
// To log out a user, clear the auth_token cookie
func Logout(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,         // Ensure the cookie is not accessible via JavaScript
        Secure:   true,         // Use "Secure" only in production with HTTPS
        SameSite: http.SameSiteStrictMode,
        MaxAge:   -1,           // Expire immediately
    })
}