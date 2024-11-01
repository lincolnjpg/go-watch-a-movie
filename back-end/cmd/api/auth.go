package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type auth struct {
	issuer        string
	audience      string
	secret        string
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
	cookie
}

type cookie struct {
	domain string
	path   string
	name   string
}

type jwtUser struct {
	id        int
	firstName string
	lastName  string
}

type tokenPairs struct {
	AccesToken   string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type claims struct {
	jwt.RegisteredClaims
}

func (j *auth) generateTokenPair(user *jwtUser) (tokenPairs, error) {
	// create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.firstName, user.lastName)
	claims["sub"] = fmt.Sprint(user.id)
	claims["aud"] = j.audience
	claims["iss"] = j.issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"

	// set the expiry for JWT
	claims["exp"] = time.Now().UTC().Add(j.tokenExpiry).Unix()

	// create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return tokenPairs{}, err
	}

	// create a refresh token and set claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.id)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()
	refreshTokenClaims["typ"] = "JWT"

	// set the expiry for refresh tokens
	claims["exp"] = time.Now().UTC().Add(j.refreshExpiry).Unix()

	// create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.secret))
	if err != nil {
		return tokenPairs{}, err
	}

	// create tokenPais and populate with signed tokens
	var tokenPairs = tokenPairs{
		AccesToken:   signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// return tokenPairs
	return tokenPairs, nil
}

func (j *auth) getRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.name,
		Path:     j.cookie.path,
		Value:    refreshToken,
		Expires:  time.Now().Add(j.refreshExpiry),
		MaxAge:   int(j.refreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   j.cookie.domain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (j *auth) getExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.name,
		Path:     j.cookie.path,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   j.cookie.domain,
		HttpOnly: true,
		Secure:   true,
	}
}

// refactor in two functions later
func (j *auth) getAndValidateJwtTokenFromHeader(w http.ResponseWriter, r *http.Request) (string, *claims, error) {
	w.Header().Add("Vary", "Authorization")

	// get auth header

	authHeader := r.Header.Get("Authorization")

	// sanity check
	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	// split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}

	// check to see if we have the word Bearer
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("invalid auth header")
	}

	token := headerParts[1]

	// declare an empty claims
	claims := &claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(j.secret), nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}

		return "", nil, err
	}

	if claims.Issuer != j.issuer {
		return "", nil, errors.New("invalid issuer")
	}

	return token, claims, nil
}
