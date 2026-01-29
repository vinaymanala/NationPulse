package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nationpulse-bff/internal/auth"
	"github.com/nationpulse-bff/internal/utils"
	"go.uber.org/zap"
)

type Middleware func(*utils.Configs, http.Handler) http.Handler
type WithAuthMiddleware func(*utils.Configs, http.Handler) http.Handler
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func checkPermissions(configs *utils.Configs, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Form.Get("userID")
		cacheID := "utils:modulePermissions:" + userID
		log.Println("ID", userID)
		id, err := strconv.Atoi(userID)
		if err != nil {
			log.Fatal("Error converting userId to int", err)
			return
		}
		var permissions []utils.UserPermissions
		data, err := utils.GetModulePermissionsFromCache(configs, id, cacheID, permissions, w, r)
		if err != nil {
			log.Println("Error getting module permissions", err)
		}
		log.Println(data)
		// check for req URL path and the corresponding permissions
		requestPath := r.URL.Path

		if !utils.HasPermissions(requestPath, &data) {
			log.Println("Error: Not Authorized to request this resources: " + requestPath)
			http.Error(w, "Forbidden resource", http.StatusForbidden)
			return
		}
		log.Println("Permission check OK! : " + requestPath)
		next.ServeHTTP(w, r)
	})
}

func allowCors(configs *utils.Configs, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin) // echo exact origin
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		log.Default().Println("CORS middleware executed for origin:", origin)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func logging(configs *utils.Configs, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		path := r.URL.Path
		query := r.URL.RawQuery

		next.ServeHTTP(w, r)

		latency := time.Since(start)
		status := &responseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}
		method := r.Method
		userAgent := r.UserAgent()

		configs.Logger.Info("request completed",
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status.statusCode),
			zap.Duration("latency", latency),
			zap.String("method", method),
			zap.String("user-agent", userAgent),
		)
	})
}

func metrics(configs *utils.Configs, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		path := r.URL.Path

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		next.ServeHTTP(w, r)

		status := fmt.Sprintf("%d", wrapped.statusCode)
		duration := time.Since(start).Seconds()
		fmt.Println("===============metrics====================", status, duration)
		configs.MetricHttpRequests.WithLabelValues(r.Method, path, status).Inc()
		configs.MetricHttpDurations.WithLabelValues(r.Method, path).Observe(duration)
	})
}

func panicRecovery(configs *utils.Configs, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func bearerFromHandler(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}

func authMiddleware(configs *utils.Configs, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement authentication logic here
		var token string
		if c, err := r.Cookie("access_token"); err == nil && c != nil && c.Value != "" {
			// fmt.Println("--------COOKIE----------", c.Value)
			token = c.Value
		} else {
			fmt.Println("No cookie found, checking Authorization header")
			token = bearerFromHandler(r)
		}
		// fmt.Println("--------TOKEN----------", token)
		if token == "" {
			log.Println(http.StatusUnauthorized, "User does not have access or missing token")
			http.Error(w, "User does not have access or missing token", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ParseAccess(token, configs)
		if err != nil {
			log.Println(http.StatusUnauthorized, err, "User does not have access")
			http.Error(w, "User does not have access", http.StatusUnauthorized)
			return
		}

		ctx := context.Background()
		if _, err := configs.Cache.GetUserByJTI(ctx, "access:"+claims.ID); err != nil {
			log.Println(http.StatusUnauthorized, err, "User does not have access or invalid token jti")
			http.Error(w, "User does not have acces or invalid token jti", http.StatusUnauthorized)
			return
		}
		// fmt.Println("")
		// fmt.Printf("CLAIMS: %v+\n", claims)
		// fmt.Println("PATH:", r.URL.Path)
		r.ParseForm()
		r.Form.Set("userID", claims.Subject)
		next.ServeHTTP(w, r)
	})
}

func MustCookie(r *http.Request, name string) (string, error) {
	val, err := r.Cookie(name)
	if err != nil || val == nil {
		return "", errors.New("missing cookie: " + name)
	}
	return val.Value, nil
}

func executeMiddlewares(configs *utils.Configs, ms []Middleware, next http.Handler) http.Handler {
	h := next
	for _, m := range ms {
		h = m(configs, h)
	}
	return h
}

func DefaultMiddlewares(configs *utils.Configs, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middlewares := []Middleware{
			allowCors,
			metrics,
			logging,
			panicRecovery,
		}
		h := executeMiddlewares(configs, middlewares, next)
		h.ServeHTTP(w, r)
	})
}

func WithAuthMiddlewares(configs *utils.Configs, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Compose auth as an inner middleware and keep the default middlewares (CORS, logging, panic recovery)
		// as the outer layer so CORS headers are always set (even when auth fails).
		middlewares := []Middleware{
			checkPermissions,
			// func(configs *utils.Configs, next http.Handler) http.Handler {
			authMiddleware,
			// },
			DefaultMiddlewares,
		}
		h := executeMiddlewares(configs, middlewares, next)
		h.ServeHTTP(w, r)
	})
}
