package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
)

func TestVersions(t *testing.T) {
	recorder := httptest.NewRecorder()
	(&API{}).Versions(recorder, httptest.NewRequest(http.MethodGet, "/api/compose/versions", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got, want := recorder.Body.String(), "/api/compose/v0"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
	if got, want := recorder.Header().Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Fatalf("Content-Type = %q, want %q", got, want)
	}
}

func TestCheckAuth(t *testing.T) {
	e := &API{Config: &cfg.Config{Credentials: cfg.UserCredentials{Username: "user", Password: "secret"}}}

	tests := []struct {
		name       string
		request    *http.Request
		wantStatus int
		wantOK     bool
	}{
		{"valid credentials", httptest.NewRequest(http.MethodGet, "/", nil), http.StatusOK, true},
		{"invalid credentials", httptest.NewRequest(http.MethodGet, "/", nil), http.StatusUnauthorized, false},
		{"missing credentials", httptest.NewRequest(http.MethodGet, "/", nil), http.StatusUnauthorized, false},
	}
	tests[0].request.SetBasicAuth("user", "secret")
	tests[1].request.SetBasicAuth("user", "wrong")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			if got := e.CheckAuth(tt.request, recorder); got != tt.wantOK {
				t.Fatalf("CheckAuth() = %v, want %v", got, tt.wantOK)
			}
			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
		})
	}
}

func TestCheckAuthAllowsUnconfiguredCredentials(t *testing.T) {
	e := &API{Config: &cfg.Config{}}
	recorder := httptest.NewRecorder()

	if !e.CheckAuth(httptest.NewRequest(http.MethodGet, "/", nil), recorder) {
		t.Fatal("CheckAuth() = false, want true when credentials are unconfigured")
	}
}

func TestIncreaseTaskCount(t *testing.T) {
	tests := map[string]string{
		"project.service.4": "project.service.5",
		"project.service":   "project.service.0",
		"":                  ".0",
	}
	for input, want := range tests {
		if got := (&API{}).IncreaseTaskCount(input); got != want {
			t.Errorf("IncreaseTaskCount(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestCommandsRoutesVersions(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/compose/versions", nil)
	(&API{}).Commands().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}
