package encryptcredentials

import (
	"app/internal/model"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// type LoginData struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

func Encrypt(loginData model.VerifyPasswordLoginInput) (string, error) {
	encryptionKey := os.Getenv("AES_256_ENCRYPTION_KEY")
	if encryptionKey == "" {
		return "", fmt.Errorf("AES_256_ENCRYPTION_KEY environment variable not set")
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize) // 16 bytes = 128 bits (required for AES block size)
	_, err = rand.Read(iv)
	if err != nil {
		return "", err
	}

	// Decode the key using base64 URL encoding (no padding)
	key, err := base64.RawURLEncoding.DecodeString(encryptionKey)
	if err != nil {
		panic("invalid base64url key: " + err.Error())
	}
	if len(key) != 32 {
		panic("key must be 32 bytes for AES-256")
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create CBC encrypter
	encrypter := cipher.NewCBCEncrypter(block, iv)

	// Pad plaintext to match block size
	paddedData := pkcs7Pad(jsonData, 16)

	// Allocate space for ciphertext
	ciphertext := make([]byte, len(paddedData))

	// Create CBC encrypter and encrypt
	encrypter.CryptBlocks(ciphertext, paddedData)

	// Final output (IV + ciphertext, base64-encoded)
	ivEncoded := base64.RawURLEncoding.EncodeToString(iv)
	ciphertextEncoded := base64.RawURLEncoding.EncodeToString(ciphertext)

	// Return the IV and ciphertext combined with a dot
	return ivEncoded + "." + ciphertextEncoded, nil
}

func Decrypt(encryptedData string) (model.VerifyPasswordLoginInput, error) {
	encryptionKey := os.Getenv("AES_256_ENCRYPTION_KEY")
	if encryptionKey == "" {
		return model.VerifyPasswordLoginInput{}, fmt.Errorf("AES_256_ENCRYPTION_KEY environment variable not set")
	}

	// Split the encrypted string on "." to extract the iv & cipher text
	parts := strings.Split(encryptedData, ".")
	if len(parts) != 2 {
		return model.VerifyPasswordLoginInput{}, fmt.Errorf("invalid encrypted data format")
	}

	// Decode both the iv and cipher text from base64 to binary
	iv, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return model.VerifyPasswordLoginInput{}, fmt.Errorf("failed to decode IV: %v", err)
	}

	ciphertext, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return model.VerifyPasswordLoginInput{}, fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	// Decode the encryption key from base64
	key, err := base64.RawURLEncoding.DecodeString(encryptionKey)
	if err != nil {
		return model.VerifyPasswordLoginInput{}, fmt.Errorf("invalid base64url encryption key: %v", err)
	}
	if len(key) != 32 {
		return model.VerifyPasswordLoginInput{}, fmt.Errorf("key must be 32 bytes for AES-256")
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return model.VerifyPasswordLoginInput{}, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Create CBC Decrypter
	decrypter := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the cipher text
	decrypted := make([]byte, len(ciphertext))
	decrypter.CryptBlocks(decrypted, ciphertext)

	// Unpad the decrypted data
	decryptedData, err := pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return model.VerifyPasswordLoginInput{}, fmt.Errorf("failed to unpad data: %v", err)
	}

	// Unmarshal the decrypted JSON into the original struct
	var loginData model.VerifyPasswordLoginInput
	err = json.Unmarshal(decryptedData, &loginData)
	if err != nil {
		return model.VerifyPasswordLoginInput{}, fmt.Errorf("failed to unmarshal decrypted data: %v", err)
	}

	return loginData, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding) // Repeat byte 'padding' for the padding length
	return append(data, padText...)                         // Append the padding to the data
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	paddingLen := int(data[len(data)-1])
	if paddingLen > blockSize || paddingLen == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:len(data)-paddingLen], nil
}
