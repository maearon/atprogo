package auth

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Claims represents the claims in a JWT
type Claims struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp"`
	NotBefore int64  `json:"nbf,omitempty"`
	IssuedAt int64  `json:"iat"`
	JWTID     string `json:"jti,omitempty"`
}

// JWT represents a JSON Web Token
type JWT struct {
	Header    map[string]interface{}
	Claims    Claims
	Signature []byte
}

// NewJWT creates a new JWT
func NewJWT(issuer, subject string, expiresIn time.Duration) *JWT {
	now := time.Now()
	return &JWT{
		Header: map[string]interface{}{
			"alg": "EdDSA",
			"typ": "JWT",
		},
		Claims: Claims{
			Issuer:    issuer,
			Subject:   subject,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(expiresIn).Unix(),
		},
	}
}

// Sign signs the JWT with a private key
func (j *JWT) Sign(privateKey ed25519.PrivateKey) (string, error) {
	// Encode header
	headerJSON, err := json.Marshal(j.Header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Encode claims
	claimsJSON, err := json.Marshal(j.Claims)
	if err != nil {
		return "", err
	}
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Create signature input
	signatureInput := headerEncoded + "." + claimsEncoded

	// Sign
	signature := ed25519.Sign(privateKey, []byte(signatureInput))
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	// Create JWT
	return headerEncoded + "." + claimsEncoded + "." + signatureEncoded, nil
}

// Verify verifies the JWT with a public key
func VerifyJWT(tokenString string, publicKey ed25519.PublicKey) (*Claims, error) {
	// Split token
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode header
	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid header encoding: %v", err)
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, fmt.Errorf("invalid header: %v", err)
	}

	// Check algorithm
	alg, ok := header["alg"].(string)
	if !ok || alg != "EdDSA" {
		return nil, fmt.Errorf("unsupported algorithm: %v", alg)
	}

	// Decode claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid claims encoding: %v", err)
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("invalid claims: %v", err)
	}

	// Decode signature
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid signature encoding: %v", err)
	}

	// Verify signature
	signatureInput := parts[0] + "." + parts[1]
	if !ed25519.Verify(publicKey, []byte(signatureInput), signature) {
		return nil, fmt.Errorf("invalid signature")
	}

	// Verify expiration
	now := time.Now().Unix()
	if claims.ExpiresAt < now {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}
