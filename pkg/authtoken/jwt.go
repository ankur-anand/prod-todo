package authtoken

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

// KID hint indicating which specific key owned by the
// signer should be used to validate the signature
// and help in validating the header.
type KID int8

const (
	// RSA256 indicates signing is done using SHA-256 hash algorithm
	RSA256 float64 = iota
)

// Claims defines custom claims that will be encoded to a JWT.
type Claims struct {
	UserID string `json:"user"` // uuid that represents user
	jwt.StandardClaims
}

// JWT Provide a JSON Web Token and Validation
type JWT struct {
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
	issuer        string
	aud           []string
	validator     *jwt.ValidationHelper
	validDuration time.Duration
}

// NewJWT return an initialized JWT
func NewJWT(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, iss, aud string, validDuration time.Duration) (JWT, error) {
	var j JWT
	if privKey != nil && pubKey != nil {
		sAud := []string{aud}
		j.validDuration = validDuration
		j.rsaPublicKey = pubKey
		j.rsaPrivateKey = privKey
		j.issuer = iss
		j.aud = sAud
		j.validator = jwt.NewValidationHelper(jwt.WithAudience(aud), jwt.WithIssuer(iss))
		return j, nil
	}
	return j, fmt.Errorf("rsa private key and public key should not be nil")
}

// Validate validates the provided token
func (j JWT) Validate(token string) (string, error) {
	c := &Claims{}

	parsedT, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		alg := token.Header["alg"].(string)
		kid := token.Header["kid"].(float64)
		if alg == "RS256" && kid == 0 {
			return j.rsaPublicKey, nil
		}
		return nil, fmt.Errorf("unexpected signing method: %s, %f", alg, kid)
	}, jwt.WithAudience(j.aud[0]))

	if err != nil {
		return "", err
	}

	// check if the token is valid for "exp, iat, nbf"
	if err := parsedT.Claims.Valid(j.validator); err != nil {
		return "", err
	}

	return c.UserID, nil
}

// Generate a new token
func (j JWT) Generate(id string) (string, error) {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(j.validDuration * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claim := &Claims{}
	claim.UserID = id
	claim.Audience = j.aud
	claim.Issuer = j.issuer
	claim.ExpiresAt = jwt.At(expirationTime)

	alg := jwt.GetSigningMethod("RS256")
	token := jwt.New(alg)
	token.Claims = claim
	token.Header = map[string]interface{}{
		"typ": "JWT",
		"alg": token.Method.Alg(),
		"kid": RSA256,
	}

	return token.SignedString(j.rsaPrivateKey)
}
