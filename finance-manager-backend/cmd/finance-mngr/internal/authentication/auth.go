package authentication

import (
	"errors"
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
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
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
	fmt.Printf("[Enter %s]\n", method)

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
		fmt.Printf("[%s] Auth header is invalid\n", method)
		fmt.Printf("[Exit %s]\n", method)
		return "", nil, errors.New("invalid auth header")
	}

	// check to see if we have the word Bearer
	if headerParts[0] != "Bearer" {
		fmt.Printf("[%s] Auth header is invalid\n", method)
		fmt.Printf("[Exit %s]\n", method)
		return "", nil, errors.New("invalid auth header")
	}

	token := headerParts[1]

	// decare an empty claims
	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("[%s] Auth header contained unexpected signing method\n", method)
			fmt.Printf("[Exit %s]\n", method)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			fmt.Printf("[%s] token is expired\n", method)
			fmt.Printf("[Exit %s]\n", method)
			return "", nil, errors.New("expired token")
		}

		fmt.Printf("[%s] unexpected error\n", method)
		fmt.Printf("[Exit %s]\n", method)
		return "", nil, err
	}

	if claims.Issuer != j.Issuer {
		fmt.Printf("[%s] invalid header issuer\n", method)
		fmt.Printf("[Exit %s]\n", method)
		return "", nil, errors.New("invalid issuer")
	}

	fmt.Printf("[Exit %s]\n", method)
	return token, claims, nil
}

func (j *Auth) GetLoggedInUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	method := "auth.GetLoggedInUserId"
	fmt.Printf("[Enter %s]\n", method)

	//This also verifies that the token is valid
	_, claims, err := j.GetTokenFromHeaderAndVerify(w, r)

	if err != nil {
		fmt.Printf("[%s] caught error when attempting to fetch claims", method)
		fmt.Printf("[Exit %s]\n", method)
		return -1, err
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		fmt.Printf("[%s] caught error when attempting to decode claims subject", method)
		fmt.Printf("[Exit %s]\n", method)
		return -1, err
	}

	fmt.Printf("[Exit %s]\n", method)
	return id, nil
}
