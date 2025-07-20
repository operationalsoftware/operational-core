package repository

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) GeneratePutTempURL(filePath string, expiresInSeconds int64, tempURLKey string) string {
	method := "PUT"
	expires := time.Now().Unix() + expiresInSeconds

	stringToSign := fmt.Sprintf("%s\n%d\n%s", method, expires, filePath)
	mac := hmac.New(sha1.New, []byte(tempURLKey))
	mac.Write([]byte(stringToSign))
	sig := hex.EncodeToString(mac.Sum(nil))

	tempURL := fmt.Sprintf("https://orbit.brightbox.com%s?temp_url_sig=%s&temp_url_expires=%d", filePath, sig, expires)
	return tempURL
}

func (r *FileRepository) GetFileURL(ctx context.Context, fileUUID uuid.UUID, expiresInSeconds int64) (string, error) {
	// expires := time.Now().Unix() + expiresInSeconds

	// fileURL must be the path, e.g., "/v1/AUTH_account/container/file.pdf"
	// stringToSign := fmt.Sprintf("%s\n%d\n%s", method, expires, fileURL)

	// mac := hmac.New(sha1.New, []byte(secretKey))
	// mac.Write([]byte(stringToSign))
	// signature := hex.EncodeToString(mac.Sum(nil))

	// tempURL := fmt.Sprintf("https://orbit.brightbox.com%s?temp_url_sig=%s&temp_url_expires=%d", fileURL, signature, expires)
	// return tempURL, nil

	return "", nil

}
