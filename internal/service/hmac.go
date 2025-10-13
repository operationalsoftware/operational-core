package service

// HMACService provides access to the application HMAC secret.
// It can be expanded later to provide helper methods for signing and validation.
type HMACService struct {
	secret string
}

// NewHMACService constructs a new HMACService with the provided secret.
func NewHMACService(secret string) *HMACService {
	return &HMACService{secret: secret}
}

// Secret returns the configured HMAC secret.
func (s *HMACService) Secret() string {
	return s.secret
}
