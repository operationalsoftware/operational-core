package apphmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Claims struct {
	Entity            string
	EntityID          string
	AllowedOperations []string
	Expires           int64
}

func GenerateHMAC(payload Claims, secret string) string {

	sort.Strings(payload.AllowedOperations)

	message := generateMessage(payload.Entity, payload.EntityID, strings.Join(payload.AllowedOperations, ","), payload.Expires)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))

	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyHMAC(claims Claims, providedHMAC, secret string) bool {
	if time.Now().Unix() > claims.Expires {
		return false
	}

	sort.Strings(claims.AllowedOperations)

	message := generateMessage(claims.Entity, claims.EntityID, strings.Join(claims.AllowedOperations, ","), claims.Expires)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expectedHMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(providedHMAC), []byte(expectedHMAC))
}

func generateMessage(entity, entityID, allowedOperations string, expiry int64) string {
	return fmt.Sprintf(
		"%s|%s|%v|%d",
		entity,
		entityID,
		allowedOperations,
		expiry)

}
