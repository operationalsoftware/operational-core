package apphmac

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

// Payload represents the data to be signed and later verified.
// It centralizes authorization context for client-to-server requests.
//
// Examples:
//   - Add a comment to thread 42:
//     Entity: "comment", EntityID: "42", Permission: "add"
//   - Add an attachment to comment 1337:
//     Entity: "comment", EntityID: "1337", Permission: "add"
//   - Notes can use: Entity: "notes", EntityID: "<notes-resource-id>", Permission: "add"
type Payload struct {
	Entity      string   `json:"Entity"`
	EntityID    string   `json:"EntityID"`
	Permissions []string `json:"Permissions"`
	Expires     int64    `json:"Expires"`
}

// Envelope contains the payload alongside its HMAC signature.
// Clients send this structure to the server, and the server verifies it.
type Envelope struct {
	Payload   Payload `json:"Payload"`
	Signature string  `json:"Signature"`
}

// SignEnvelope produces an Envelope by signing the given payload with the provided secret.
// It uses the existing HMAC scheme to generate the signature.
func SignEnvelope(p Payload, secret string) Envelope {
	claims := Payload{
		Entity:      p.Entity,
		EntityID:    p.EntityID,
		Permissions: p.Permissions,
		Expires:     p.Expires,
	}
	sig := GenerateHMAC(claims, secret)
	return Envelope{Payload: p, Signature: sig}
}

// VerifyEnvelope verifies the signature for the provided envelope against the secret.
// It returns the payload if valid, otherwise an error.
func VerifyEnvelope(e Envelope, secret string) (Payload, error) {
	if e.Signature == "" {
		return Payload{}, errors.New("missing signature")
	}
	claims := Payload{
		Entity:      e.Payload.Entity,
		EntityID:    e.Payload.EntityID,
		Permissions: e.Payload.Permissions,
		Expires:     e.Payload.Expires,
	}
	if ok := VerifyHMAC(claims, e.Signature, secret); !ok {
		return Payload{}, errors.New("invalid signature")
	}
	return e.Payload, nil
}

// VerifySignature verifies a signature against a payload without constructing an Envelope.
func VerifySignature(p Payload, signature, secret string) bool {
	claims := Payload{
		Entity:      p.Entity,
		EntityID:    p.EntityID,
		Permissions: p.Permissions,
		Expires:     p.Expires,
	}
	return VerifyHMAC(claims, signature, secret)
}

// Checks if the envelope is valid and contains the required entity, entity ID, and permission.
func IsValidEnvelope(e Envelope, secret, requiredEntity, requiredEntityID, requiredPermission string) (bool, error) {
	// Check if envelope is nil or signature is empty
	if e.Signature == "" {
		return false, fmt.Errorf("missing HMAC")
	}
	payload, err := VerifyEnvelope(e, secret)
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
