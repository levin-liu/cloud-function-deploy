package visit_count

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
)

func TestVisitCount(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run: %v", err)
	}
	defer s.Close()

	os.Setenv("REDISHOST", s.Host())
	os.Setenv("REDISPORT", s.Port())

	req := httptest.NewRequest("GET", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()

	visitCount(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("VisitCount got status %v, want %v", rr.Code, http.StatusOK)
	}
}
