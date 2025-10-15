package apphmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"
)

// Helper functions to generate and verify HMAC signatures
func generateHMAC(payload Payload, secret string) string {

	sort.Strings(payload.Permissions)

	message := generateMessage(payload.Entity, payload.EntityID, strings.Join(payload.Permissions, ","), payload.Expires)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))

	return hex.EncodeToString(mac.Sum(nil))
}

func verifyHMAC(payload Payload, providedHMAC, secret string) bool {
	if time.Now().Unix() > payload.Expires {
		return false
	}

	sort.Strings(payload.Permissions)

	message := generateMessage(payload.Entity, payload.EntityID, strings.Join(payload.Permissions, ","), payload.Expires)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expectedHMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(providedHMAC), []byte(expectedHMAC))
}

func generateMessage(entity, entityID, permissions string, expiry int64) string {
	return fmt.Sprintf(
		"%s|%s|%v|%d",
		entity,
		entityID,
		permissions,
		expiry)

}

func verifyEnvelope(e envelope, secret string) (Payload, error) {
	if e.Signature == "" {
		return Payload{}, errors.New("missing signature")
	}
	claims := Payload{
		Entity:      e.Payload.Entity,
		EntityID:    e.Payload.EntityID,
		Permissions: e.Payload.Permissions,
		Expires:     e.Payload.Expires,
	}
	if ok := verifyHMAC(claims, e.Signature, secret); !ok {
		return Payload{}, errors.New("invalid signature")
	}
	return e.Payload, nil
}

// Payload represents the data to be signed and later verified.
// It centralizes authorization context for client-to-server requests.
//
// Examples:
//   - Add a comment to thread 42:
//     entity: "comment", entity_id: "42", permissions: "add"
//   - Add an attachment to comment 1337:
//     entity: "comment", entity_id: "1337", permissions: "add"
//   - Notes can use: entity: "notes", entity_id: "<notes-resource-id>", permissions: "add"
type Payload struct {
	Entity      string   `json:"entity"`
	EntityID    string   `json:"entity_id"`
	Permissions []string `json:"permissions"`
	Expires     int64    `json:"expires"`
}

// Envelope contains the payload alongside its HMAC signature.
// Clients send this structure to the server, and the server verifies it.
type envelope struct {
	Payload   Payload `json:"payload"`
	Signature string  `json:"signature"`
}

// AppHMAC provides access to the HMAC functionality.
type AppHMAC struct {
	secret string
}

// NewAppHMAC constructs a new HMACService with the provided secret.
func NewAppHMAC(secret string) *AppHMAC {
	return &AppHMAC{secret: secret}
}

func (s *AppHMAC) CreateEnvelope(p Payload) string {
	claims := Payload{
		Entity:      p.Entity,
		EntityID:    p.EntityID,
		Permissions: p.Permissions,
		Expires:     p.Expires,
	}
	sig := generateHMAC(claims, s.secret)
	env := envelope{Payload: p, Signature: sig}
	jsonEnv, _ := json.Marshal(env)
	return string(jsonEnv)
}

// extract permissions from envelope
func (s *AppHMAC) GetEnvelopePermissions(envelopeStr string) ([]string, error) {
	var e envelope
	if err := json.Unmarshal([]byte(envelopeStr), &e); err != nil {
		return nil, fmt.Errorf("invalid envelope json: %w", err)
	}
	if e.Signature == "" {
		return nil, fmt.Errorf("missing HMAC")
	}
	claims, err := verifyEnvelope(e, s.secret)
	if err != nil {
		return nil, fmt.Errorf("error validating HMAC: %w", err)
	}
	return claims.Permissions, nil
}

func (s *AppHMAC) CheckEnvelope(envelopeStr string, requiredEntity, requiredEntityID, requiredPermission string) (bool, error) {
	// Check if envelope is nil or signature is empty
	// decode envelope JSON
	var e envelope
	if err := json.Unmarshal([]byte(envelopeStr), &e); err != nil {
		return false, fmt.Errorf("invalid envelope json: %w", err)
	}
	if e.Signature == "" {
		return false, fmt.Errorf("missing HMAC")
	}
	payload, err := verifyEnvelope(e, s.secret)
	// Check if verification failed
	if err != nil {
		return false, fmt.Errorf("error validating HMAC: %w", err)
	}
	// Check if payload matches required entity, entity ID
	if payload.Entity != requiredEntity || payload.EntityID != requiredEntityID ||
		!slices.Contains(payload.Permissions, requiredPermission) {
		return false, errors.New("envelope does not grant required permission")
	}

	// Check if payload is expired
	if payload.Expires < (time.Now().Unix()) {
		return false, errors.New("expired")
	}

	// All checks passed
	return true, nil
}
