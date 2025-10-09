package apphmac

import (
	"errors"
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
	Entity      string   `json:"entity"`
	EntityID    string   `json:"entityId"`
	Permissions []string `json:"permissions"`
	Expires     int64    `json:"expires"`
}

// Envelope contains the payload alongside its HMAC signature.
// Clients send this structure to the server, and the server verifies it.
type Envelope struct {
	Payload   Payload `json:"payload"`
	Signature string  `json:"signature"`
}

// SignEnvelope produces an Envelope by signing the given payload with the provided secret.
// It uses the existing HMAC scheme by mapping the single Permission into AllowedOperations.
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
