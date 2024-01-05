package authentication

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type JwtUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Roles     string `json:"roles"`
}

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
}

func (j *Auth) GenerateTokenPair(user *JwtUser) (TokenPairs, error) {
	// Create a Token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s, %s", user.LastName, user.FirstName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"
	claims["roles"] = fmt.Sprint(user.Roles)

	// Set the expiry for JWT
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create a Refresh token an dset claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()

	// Set the expiry for the refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	// Create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create TokenPairs and populate with signed tokens
	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// Return token pairs
	return tokenPairs, nil
}

func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    refreshToken,
		Expires:  time.Now().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (j *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	method := "auth.GetTokenFromHeaderAndVerify"
	fmlogger.Enter(method)

	w.Header().Add("Vary", "Authorization")

	// Get Auth Header
	authHeader := r.Header.Get("Authorization")

	// Sanity Check
	if authHeader == "" {
		fmt.Printf("[%s] Auth header is not present\n", method)
		fmt.Printf("[Exit %s]\n", method)
		return "", nil, errors.New("no auth header")
	}

	// Split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		err := errors.New(constants.InvalidAuthHeaderError)
		fmlogger.ExitError(method, err.Error(), err)
		return "", nil, err
	}

	// check to see if we have the word Bearer
	if headerParts[0] != "Bearer" {
		err := errors.New(constants.InvalidAuthHeaderError)
		fmlogger.ExitError(method, err.Error(), err)
		return "", nil, err
	}

	token := headerParts[1]
	claims := &Claims{}

	token, claims, err := j.ParseAndVerifyToken(token)

	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		return "", nil, err
	}

	fmlogger.Exit(method)
	return token, claims, nil

}

func (j *Auth) ParseAndVerifyToken(token string) (string, *Claims, error) {
	method := "auth.ParseAndVerifyToken"
	fmlogger.Enter(method)

	// decare an empty claims
	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmlogger.ExitError(method, constants.InvalidSigningMethodError, errors.New(constants.InvalidSigningMethodError))
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			fmlogger.ExitError(method, constants.ExpiredTokenError, err)
			return "", nil, errors.New(constants.ExpiredTokenError)
		}

		fmlogger.ExitError(method, "unexpected error", err)
		return "", nil, err
	}

	if claims.Issuer != j.Issuer {
		err := errors.New(constants.InvalidIssuerError)
		fmlogger.ExitError(method, err.Error(), err)
		return "", nil, err
	}

	fmlogger.Exit(method)
	return token, claims, nil
}

func (j *Auth) GetLoggedInUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	method := "auth.GetLoggedInUserId"
	fmlogger.Enter(method)

	//This also verifies that the token is valid
	_, claims, err := j.GetTokenFromHeaderAndVerify(w, r)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error fetching claims", err)
		return -1, err
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		fmlogger.ExitError(method, "unexpected error decoding claims subject", err)
		return -1, err
	}

	fmlogger.Exit(method)
	return id, nil
}
