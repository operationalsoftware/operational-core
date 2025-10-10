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

func GenerateHMAC(payload Payload, secret string) string {

	sort.Strings(payload.Permissions)

	message := generateMessage(payload.Entity, payload.EntityID, strings.Join(payload.Permissions, ","), payload.Expires)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))

	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyHMAC(payload Payload, providedHMAC, secret string) bool {
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
