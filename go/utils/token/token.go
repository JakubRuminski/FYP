package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

func ValidToken(logger *logger.Logger, r *http.Request) (ok bool) {
	// Read the token from the Auth-Token cookie
	cookie, err := r.Cookie("Auth-Token")
	if err != nil {
		logger.ERROR("Failed to get cookie. Reason: %s", err)
		return false
	}

	if cookie.Value == "" {
		logger.ERROR("Cookie value is empty")
		return false
	}

	tokenParts := strings.Split(cookie.Value, ".")
	if len(tokenParts) != 2 {
		logger.ERROR("Invalid token structure")
		return false
	}

	jsonToken, err := base64.StdEncoding.DecodeString(tokenParts[0])
	if err != nil {
		logger.ERROR("Failed to decode token. Reason: %s", err)
		return false
	}

	token := &Token{}
	err = json.Unmarshal(jsonToken, token)
	if err != nil {
		logger.ERROR("Failed to unmarshal token. Reason: %s", err)
		return false
	}

	// Verify the HMAC
	tokenKey, ok := env.Get(logger, "TOKEN_KEY")
	if !ok {
		logger.ERROR("Failed to get token key")
		return false
	}

	h := hmac.New(sha256.New, []byte(tokenKey))
	h.Write(jsonToken)
	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if sha != tokenParts[1] {
		logger.ERROR("Invalid token signature")
		return false
	}

	// Check the token expiration time
	if time.Unix(token.ExpiresAt, 0).Before(time.Now()) {
		logger.ERROR("Token expired")
		return false
	}

	// Token is valid
	return true
}

func CreateToken(logger *logger.Logger, w http.ResponseWriter, user string) (ok bool) {

	ok = createToken(logger, w, user)
	if !ok {
		logger.ERROR("Failed to create token")
		return false
	}

	return true
}

func createToken(logger *logger.Logger, w http.ResponseWriter, user string) (ok bool) {

	token := &Token{
		User:      user,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	jsonToken, err := json.Marshal(token)
	if err != nil {
		logger.ERROR("Failed to marshal token. Reason: %s", err)
		return false
	}

	tokenKey, ok := env.Get(logger, "TOKEN_KEY")
	if !ok {
		logger.ERROR("Failed to get token key")
		return false
	}

	h := hmac.New(sha256.New, []byte(tokenKey))

	h.Write(jsonToken)

	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))

	encoding := base64.StdEncoding.EncodeToString(jsonToken) + "." + sha

	// Create a new cookie with the token
	cookie := &http.Cookie{
		Name:     "Auth-Token",
		Value:    encoding,
		HttpOnly: true,
		Secure:   true, // Set this to true if you are using HTTPS
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	}

	// Set the cookie in the HTTP response
	http.SetCookie(w, cookie)

	return true
}

type Token struct {
	User      string `json:"user"`
	ExpiresAt int64  `json:"expiresAt"`
}
