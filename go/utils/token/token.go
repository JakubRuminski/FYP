package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

func ValidToken(logger *logger.Logger, r *http.Request) (token *Token, ok bool) {
	// Read the token from the Auth-Token cookie
	cookie, err := r.Cookie("Auth-Token")
	if err != nil {
		logger.ERROR("Failed to get cookie. Reason: %s", err)
		return nil, false
	}

	if cookie.Value == "" {
		logger.ERROR("Cookie value is empty")
		return nil, false
	}

	tokenParts := strings.Split(cookie.Value, ".")
	if len(tokenParts) != 2 {
		logger.ERROR("Invalid token structure")
		return nil, false
	}

	jsonToken, err := base64.StdEncoding.DecodeString(tokenParts[0])
	if err != nil {
		logger.ERROR("Failed to decode token. Reason: %s", err)
		return nil, false
	}

	token = &Token{}
	err = json.Unmarshal(jsonToken, token)
	if err != nil {
		logger.ERROR("Failed to unmarshal token. Reason: %s", err)
		return nil, false
	}

	// Verify the HMAC
	tokenKey, ok := env.Get(logger, "TOKEN_KEY")
	if !ok {
		logger.ERROR("Failed to get token key")
		return nil, false
	}

	h := hmac.New(sha256.New, []byte(tokenKey))
	h.Write(jsonToken)
	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if sha != tokenParts[1] {
		logger.ERROR("Invalid token signature")
		return nil, false
	}

	// Check the token expiration time
	if time.Unix(token.ExpiresAt, 0).Before(time.Now()) {
		logger.ERROR("Token expired")
		return nil, false
	}

	// Token is valid
	return token, true
}

func CreateToken(logger *logger.Logger, w http.ResponseWriter) (ok bool) {

    src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	pseudorandomID := rnd.Int()

	ok = createToken(logger, w, pseudorandomID)
	if !ok {
		logger.ERROR("Failed to create token")
		return false
	}

	return true
}

func createToken(logger *logger.Logger, w http.ResponseWriter, pseudorandomID int) (ok bool) {

	TOKEN_EXPIRY, ok := env.GetInt(logger, "TOKEN_EXPIRY")
	if !ok {
		return false
	}

	token := &Token{
		PseudorandomID:  pseudorandomID,
		ExpiresAt:       time.Now().Add( time.Duration(TOKEN_EXPIRY) * time.Hour ).Unix(),
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
		Secure:   true,
		Expires:  time.Now().Add(time.Duration(TOKEN_EXPIRY) * time.Hour),
		Path:     "/",
	}

	// Set the cookie in the HTTP response
	http.SetCookie(w, cookie)

	return true
}

type Token struct {
	PseudorandomID      int      `json:"pseudorandom_id"`
	ExpiresAt           int64    `json:"expires_at"`
}


func GetID(logger *logger.Logger, r *http.Request) (int, bool) {
    token, valid := ValidToken(logger, r)
    if !valid {
        logger.ERROR("Invalid or expired token")
        return 0, false
    }

    return token.PseudorandomID, true
}