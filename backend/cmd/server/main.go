package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	validationCode = "WAITLIST_VALIDATION_ERROR"
	internalCode   = "WAITLIST_INTERNAL_ERROR"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type errorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
type request struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type waitlistStore interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}
type limiter struct {
	mu   sync.Mutex
	hits map[string][]time.Time
	now  func() time.Time
}

func newLimiter(now func() time.Time) *limiter {
	if now == nil {
		now = time.Now
	}
	return &limiter{hits: map[string][]time.Time{}, now: now}
}
func (l *limiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	kept := make([]time.Time, 0, len(l.hits[ip]))
	for _, t := range l.hits[ip] {
		if now.Sub(t) < time.Minute {
			kept = append(kept, t)
		}
	}
	if len(kept) >= 5 {
		l.hits[ip] = kept
		return false
	}
	l.hits[ip] = append(kept, now)
	return true
}
func writeError(w http.ResponseWriter, status int, code, message string) {
	var b errorBody
	b.Error.Code = code
	b.Error.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(b)
}
func validAllowedOrigin(value string) bool {
	u, err := url.ParseRequestURI(value)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" && u.Path == "" && u.RawQuery == "" && u.Fragment == ""
}
func newRouter(store waitlistStore, allowedOrigin string, l *limiter) http.Handler {
	r := chi.NewRouter()
	r.Use(cors(allowedOrigin))
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"status":"ok"}`)
	})
	r.Post("/api/waitlist", func(w http.ResponseWriter, req *http.Request) {
		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			ip = req.RemoteAddr
		}
		if !l.allow(ip) {
			writeError(w, http.StatusTooManyRequests, "WAITLIST_RATE_LIMITED", "Please wait a moment before trying again.")
			return
		}
		decoder := json.NewDecoder(http.MaxBytesReader(w, req.Body, 4096))
		decoder.DisallowUnknownFields()
		var in request
		if decoder.Decode(&in) != nil || decoder.Decode(&struct{}{}) != io.EOF {
			writeError(w, http.StatusBadRequest, validationCode, "Enter a valid email.")
			return
		}
		in.Name = strings.TrimSpace(in.Name)
		in.Email = strings.ToLower(strings.TrimSpace(in.Email))
		if len(in.Name) > 100 || len(in.Email) == 0 || len(in.Email) > 255 || !emailPattern.MatchString(in.Email) {
			writeError(w, http.StatusBadRequest, validationCode, "Enter a valid email.")
			return
		}
		// name is optional: store SQL NULL rather than an empty string when
		// the caller omits it (e.g. the email-only /waitlist page).
		var namePtr any
		if in.Name != "" {
			namePtr = in.Name
		}
		_, err = store.Exec(req.Context(), "INSERT INTO waitlist_entries (name,email_normalized) VALUES ($1,$2)", namePtr, in.Email)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				writeError(w, http.StatusConflict, "WAITLIST_EMAIL_EXISTS", "This email is already on the waitlist.")
			} else {
				writeError(w, http.StatusInternalServerError, internalCode, "We could not join the waitlist. Please try again.")
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{"status":"created"}`)
	})
	return r
}
func cors(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			origin := req.Header.Get("Origin")
			if origin != "" && origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}
			if req.Method == http.MethodOptions {
				if origin == "" || origin != allowedOrigin {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}
func main() {
	dsn, allowedOrigin := os.Getenv("DATABASE_URL"), os.Getenv("ALLOWED_ORIGIN")
	if dsn == "" {
		slog.Error("DATABASE_URL is required")
		os.Exit(1)
	}
	if !validAllowedOrigin(allowedOrigin) {
		slog.Error("ALLOWED_ORIGIN must be an absolute http(s) origin without path")
		os.Exit(1)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		slog.Error("database configuration failed")
		os.Exit(1)
	}
	defer pool.Close()
	if err = pool.Ping(context.Background()); err != nil {
		slog.Error("database unavailable")
		os.Exit(1)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("server started", "port", port)
	if err := http.ListenAndServe(":"+port, newRouter(pool, allowedOrigin, newLimiter(nil))); err != nil {
		slog.Error("server stopped", "error", err)
	}
}
