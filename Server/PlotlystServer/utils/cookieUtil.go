package utils

import (
	"net/http"
	"os"
	"strings"
)

func SetSessionCookie(writer http.ResponseWriter, name, value string, maxAge int) {
	secure := strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true")
	sameSite := http.SameSiteLaxMode
	switch strings.ToLower(strings.TrimSpace(os.Getenv("COOKIE_SAME_SITE"))) {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	}
	if sameSite == http.SameSiteNoneMode {
		secure = true
	}
	http.SetCookie(writer, &http.Cookie{Name: name, Value: value, Path: "/", Domain: strings.TrimSpace(os.Getenv("COOKIE_DOMAIN")), MaxAge: maxAge, Secure: secure, HttpOnly: true, SameSite: sameSite})
}
