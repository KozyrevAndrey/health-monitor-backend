package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func newAuthServer() *Server {
	s := &Server{
		log:        zerolog.Nop(),
		enableAuth: true,
		authUser:   "admin",
		authPass:   "secret",
		sessions:   newSessionStore(),
	}
	s.router = s.setupRouter()
	return s
}

func doLogin(t *testing.T, s *Server, body string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(rec, req)
	return rec
}

func sessionCookie(t *testing.T, rec *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, c := range rec.Result().Cookies() {
		if c.Name == sessionCookieName {
			return c
		}
	}
	return nil
}

func TestLogin_ValidCredentials_SetsSessionCookie(t *testing.T) {
	s := newAuthServer()
	rec := doLogin(t, s, `{"username":"admin","password":"secret"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", rec.Code, rec.Body.String())
	}
	c := sessionCookie(t, rec)
	if c == nil {
		t.Fatal("expected session cookie to be set")
	}
	if !c.HttpOnly {
		t.Error("session cookie must be HttpOnly")
	}
	if !s.sessions.valid(c.Value) {
		t.Error("issued session token should be valid")
	}
}

func TestLogin_InvalidCredentials_Unauthorized(t *testing.T) {
	s := newAuthServer()
	rec := doLogin(t, s, `{"username":"admin","password":"wrong"}`)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if sessionCookie(t, rec) != nil {
		t.Error("no session cookie should be set on failed login")
	}
}

func TestAuthMiddleware_APIWithoutSession_401(t *testing.T) {
	s := newAuthServer()
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/targets", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_PageWithoutSession_RedirectsToLogin(t *testing.T) {
	s := newAuthServer()
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/login" {
		t.Errorf("expected redirect to /login, got %q", loc)
	}
}

func TestAuthMiddleware_WithSession_PassesThrough(t *testing.T) {
	s := newAuthServer()
	login := doLogin(t, s, `{"username":"admin","password":"secret"}`)
	c := sessionCookie(t, login)
	if c == nil {
		t.Fatal("login did not set a session cookie")
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.AddCookie(c)
	s.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"authenticated":true`) {
		t.Errorf("expected authenticated:true, got %s", rec.Body.String())
	}
}

func TestLogout_InvalidatesSession(t *testing.T) {
	s := newAuthServer()
	login := doLogin(t, s, `{"username":"admin","password":"secret"}`)
	c := sessionCookie(t, login)
	if c == nil {
		t.Fatal("login did not set a session cookie")
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.AddCookie(c)
	s.router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on logout, got %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/targets", nil)
	req.AddCookie(c)
	s.router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 after logout, got %d", rec.Code)
	}
}

func TestLoginEndpoint_PublicWhenAuthEnabled(t *testing.T) {
	s := newAuthServer()
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/login", nil))

	// Served without a session (the handler reads the login page from disk;
	// a 200 or a 404 in the test working dir both mean the middleware let it
	// through — anything but 302/401 would do, so assert not-redirected).
	if rec.Code == http.StatusFound || rec.Code == http.StatusUnauthorized {
		t.Fatalf("login page must be reachable without a session, got %d", rec.Code)
	}
}

func TestAuthDisabled_NoLoginRequired(t *testing.T) {
	s := &Server{log: zerolog.Nop(), enableAuth: false, sessions: newSessionStore()}
	s.router = s.setupRouter()

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"auth_enabled":false`) {
		t.Errorf("expected auth_enabled:false, got %s", rec.Body.String())
	}
}
