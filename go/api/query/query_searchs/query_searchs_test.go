package query_searchs

import (
	"testing"

	"github.com/jakubruminski/FYP/go/utils/logger"
)


func TestInit(t *testing.T) {
	logger := &logger.Logger{}

	ok := INIT(logger)
	if !ok {
		t.Errorf("Failed to initialize searches")
	}
}