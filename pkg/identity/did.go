package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// DID represents a decentralized identifier
type DID struct {
	Method     string
	Identifier string
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// NewDID creates a new DID with a generated key pair
func NewDID(method string) (*DID, error) {
	// Generate key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Encode public key as base64
	pubKeyEncoded := base64.StdEncoding.EncodeToString(publicKey)

	return &DID{
		Method:     method,
		Identifier: pubKeyEncoded[:16], // Use first 16 chars of public key as identifier
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

// String returns the string representation of the DID
func (d *DID) String() string {
	return fmt.Sprintf("did:%s:%s", d.Method, d.Identifier)
}

// Sign signs data with the DID's private key
func (d *DID) Sign(data []byte) ([]byte, error) {
	if d.PrivateKey == nil {
		return nil, fmt.Errorf("private key not available")
	}
	return ed25519.Sign(d.PrivateKey, data), nil
}

// Verify verifies a signature with the DID's public key
func (d *DID) Verify(data, signature []byte) bool {
	return ed25519.Verify(d.PublicKey, data, signature)
}

// ParseDID parses a DID string
func ParseDID(didStr string) (*DID, error) {
	// In a real implementation, this would parse the DID string
	// and resolve the DID document to get the public key
	return nil, fmt.Errorf("not implemented")
}
