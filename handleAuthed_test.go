package main

import (
	"fmt"
	mis "github.com/matryer/is"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testAuthVal = 0

func testAuthHandler(w http.ResponseWriter, r *http.Request) {
	testAuthVal++
}

type testAuthVerifier struct{}

func (t testAuthVerifier) verify(idToken string) (string, error) {
	if idToken == "TestAuthPass" {
		return "user1369", nil
	}
	return "", fmt.Errorf("invalid authentication token")
}

func TestAuth(t *testing.T) {
	is := mis.New(t)

	handlerOpt := func(s *server) {
		s.handler = s.handleAuthed(testAuthVerifier{}, testAuthHandler)
	}
	srv := newServer(handlerOpt)

	tests := []struct {
		token string
		val   int
		code  int
	}{
		{"Bearer TestAuthPass", 1, http.StatusOK},
		{"Bearer ab76c0239c203c020", 1, http.StatusUnauthorized},
		{"Bearerab76c0239c203c020", 1, http.StatusUnauthorized},
		{" Bearer TestAuthPass", 1, http.StatusUnauthorized},
		{"Bearer		ab76c0239c203c020", 1, http.StatusUnauthorized},
		{"Bearer  TestAuthPass", 2, http.StatusOK},
		{"", 2, http.StatusUnauthorized},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", test.token)
		srv.ServeHTTP(rr, req)

		is.Equal(rr.Code, test.code)
		is.Equal(testAuthVal, test.val)
	}
}