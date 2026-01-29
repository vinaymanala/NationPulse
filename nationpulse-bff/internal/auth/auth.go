package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nationpulse-bff/internal/store"
	"github.com/nationpulse-bff/internal/utils"
)

type Tokens struct {
	Access   string
	Refresh  string
	JTIAcc   string
	JTIRef   string
	ExpAcc   time.Time
	ExpRef   time.Time
	UserID   string
	Issuer   string
	Audience string
}

func IssueTokens(userID string, configs *utils.Configs) (*Tokens, error) {
	now := time.Now().UTC()
	t := &Tokens{
		UserID:   userID,
		JTIAcc:   uuid.NewString(),
		JTIRef:   uuid.NewString(),
		ExpAcc:   now.Add(15 * time.Minute),
		ExpRef:   now.Add(7 * 24 * time.Hour),
		Issuer:   "jwt-nation-pulse-server",
		Audience: "jwt-nation-pulse-client",
	}

	acc := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIAcc,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{t.Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpAcc),
	})

	ref := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIRef,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{t.Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpRef),
	})

	var err error
	t.Access, err = acc.SignedString([]byte(configs.Cfg.AccessSecret))
	if err != nil {
		return nil, err
	}

	t.Refresh, err = ref.SignedString([]byte(configs.Cfg.RefreshSecret))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func Persist(ctx context.Context, r *store.Redis, t *Tokens) error {
	if err := r.SetJTI(ctx, "access:"+t.JTIAcc, t.UserID, t.ExpAcc); err != nil {
		return err
	}

	if err := r.SetJTI(ctx, "refresh:"+t.JTIRef, t.UserID, t.ExpRef); err != nil {
		return err
	}
	return nil
}

func SetAuthCookies(w http.ResponseWriter, t *Tokens) {
	fmt.Println("Cookie setting up")
	access_cookie := &http.Cookie{
		Name:     "access_token",
		Value:    t.Access,
		Path:     "/",
		Expires:  t.ExpAcc,
		HttpOnly: true,
		Secure:   false, // true in production environments
		SameSite: http.SameSiteLaxMode,
	}
	refresh_cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    t.Refresh,
		Path:     "/",
		Expires:  t.ExpRef,
		HttpOnly: true,
		Secure:   false, // true in production environments
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, access_cookie)
	http.SetCookie(w, refresh_cookie)
	fmt.Println("Cookie Set!")
}

func ClearAuthCookies(w http.ResponseWriter) {
	access_cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	refresh_cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, access_cookie)
	http.SetCookie(w, refresh_cookie)
	fmt.Println("Cookie cleared...", access_cookie, refresh_cookie)
}

func ParseAccess(tokenStr string, configs *utils.Configs) (*jwt.RegisteredClaims, error) {
	secret := configs.Cfg.AccessSecret
	// fmt.Println("SECRET RECOVERED", secret)
	return parseWithSecret(tokenStr, secret)
}

func ParseRefresh(tokenStr string, configs *utils.Configs) (*jwt.RegisteredClaims, error) {
	secret := configs.Cfg.RefreshSecret
	// fmt.Println("SECRET RECOVERED", secret)
	return parseWithSecret(tokenStr, secret)
}

func parseWithSecret(tokenStr, secret string) (*jwt.RegisteredClaims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Extra safety: ensure HMAC family
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func MustCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil || c.Value == "" {
		return "", errors.New("missing cookie" + name)
	}
	return c.Value, nil
}
