package token

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

var (
	nowTime time.Time
	logFile *os.File
	log     *logger.Logger
)

func init() {
	nowTime = time.Now()

	log = &logger.Logger{}
	var ok bool
	logFile, ok = log.InitRequestLogFile(uuid.New().String())
	if !ok {
		log.ERROR("Error while creating log file")
	}
}


func Test_ValidToken_NO_TOKEN(t *testing.T) {
	r := httptest.NewRequest("GET", "/api/search", nil)
	_, ok := ValidToken(log, r)

	if ok {
		t.Errorf("Expected false, got %t", ok)
	}

}

func Test_ValidToken_TOKEN_EXISTS(t *testing.T) {
	w := httptest.NewRecorder()
	ok := CreateToken(log, w, log.ClientID)
	if !ok {
		t.Errorf("Expected true, got %t", ok)
	}

	r := httptest.NewRequest("GET", "/api/search", nil)

	// Transfer the cookies from the response to the request
	cookies := w.Result().Cookies()
    for _, cookie := range cookies {
        r.AddCookie(cookie) 
    }

    _, ok = ValidToken(log, r)
    if !ok {
        t.Errorf("Expected true, got %t", ok)
    }

}


func Test_CreateToken(t *testing.T) {

	w := httptest.NewRecorder()

	ok := CreateToken(log, w, log.ClientID)
	if !ok {
		t.Errorf("Expected true, got %t", ok)
	}

}