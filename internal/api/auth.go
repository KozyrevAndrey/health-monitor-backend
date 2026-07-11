package api

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	sessionCookieName = "hm_session"
	sessionTTL        = 7 * 24 * time.Hour
)

// sessionStore keeps login sessions in memory. Sessions do not survive a
// restart — users simply log in again, which is acceptable for a
// single-operator dashboard and avoids persisting tokens in the database.
type sessionStore struct {
	mu       sync.Mutex
	sessions map[string]time.Time // token -> expiry
}

func newSessionStore() *sessionStore {
	return &sessionStore{sessions: make(map[string]time.Time)}
}

func (st *sessionStore) create() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	token := hex.EncodeToString(buf)

	st.mu.Lock()
	defer st.mu.Unlock()
	// Lazily drop expired sessions so the map doesn't grow unbounded.
	now := time.Now()
	for t, exp := range st.sessions {
		if now.After(exp) {
			delete(st.sessions, t)
		}
	}
	st.sessions[token] = now.Add(sessionTTL)
	return token, nil
}

func (st *sessionStore) valid(token string) bool {
	st.mu.Lock()
	defer st.mu.Unlock()
	exp, ok := st.sessions[token]
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		delete(st.sessions, token)
		return false
	}
	return true
}

func (st *sessionStore) delete(token string) {
	st.mu.Lock()
	defer st.mu.Unlock()
	delete(st.sessions, token)
}

// isPublicPath lists routes reachable without a session: the login page and
// endpoint, and static assets (the login page needs the stylesheets).
func isPublicPath(path string) bool {
	return path == "/login" ||
		path == "/api/v1/auth/login" ||
		strings.HasPrefix(path, "/static/")
}

func (s *Server) hasValidSession(r *http.Request) bool {
	c, err := r.Cookie(sessionCookieName)
	if err != nil {
		return false
	}
	return s.sessions.valid(c.Value)
}

// sessionAuthMiddleware requires a valid session cookie on every route except
// the public ones. API requests get a 401 JSON response; page requests are
// redirected to the login form.
func (s *Server) sessionAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r.URL.Path) || s.hasValidSession(r) {
			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api/") {
			s.respondError(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if !s.enableAuth {
		s.respondError(w, http.StatusBadRequest, "Authentication is disabled", nil)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	userOK := subtle.ConstantTimeCompare([]byte(req.Username), []byte(s.authUser)) == 1
	passOK := subtle.ConstantTimeCompare([]byte(req.Password), []byte(s.authPass)) == 1
	if !userOK || !passOK {
		s.log.Warn().Str("remote_addr", r.RemoteAddr).Msg("Failed login attempt")
		s.respondError(w, http.StatusUnauthorized, "Invalid username or password", nil)
		return
	}

	token, err := s.sessions.create()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to create session", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
	})
	s.respondJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookieName); err == nil {
		s.sessions.delete(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	s.respondJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// handleAuthStatus lets the frontend decide whether to show the logout button.
func (s *Server) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, http.StatusOK, map[string]bool{
		"auth_enabled":  s.enableAuth,
		"authenticated": !s.enableAuth || s.hasValidSession(r),
	})
}

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	// Nothing to log into when auth is off or the user is already signed in.
	if !s.enableAuth || s.hasValidSession(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "./web/static/login.html")
}
