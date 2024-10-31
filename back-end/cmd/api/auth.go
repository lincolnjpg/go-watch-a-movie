package main

import (
	"fmt"
	"net/http"
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
	id        int    `json:"id"`
	firstName string `json:"first_name"`
	lastName  string `json:"last_name"`
}

type tokenPairs struct {
	accesToken   string `json:"access_token"`
	refreshToken string `json:"refresh_token"`
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
		accesToken:   signedAccessToken,
		refreshToken: signedRefreshToken,
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
