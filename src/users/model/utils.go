package userModel

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		return "", fmt.Errorf("password length must be at least 8 characters")
	}

	const (
		lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
		uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers          = "0123456789"
		symbols          = "!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~"
	)

	allChars := lowercaseLetters + uppercaseLetters + numbers + symbols

	// Use time-based seed for randomness
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	password := make([]byte, length)

	// Ensure at least one lowercase, one uppercase, one number, and one symbol
	password[0] = lowercaseLetters[randomGenerator.Intn(len(lowercaseLetters))]
	password[1] = uppercaseLetters[randomGenerator.Intn(len(uppercaseLetters))]
	password[2] = numbers[randomGenerator.Intn(len(numbers))]
	password[3] = symbols[randomGenerator.Intn(len(symbols))]

	// Fill the rest of the password randomly
	for i := 4; i < length; i++ {
		password[i] = allChars[randomGenerator.Intn(len(allChars))]
	}

	// Shuffle the password to randomize the order
	randomGenerator.Shuffle(length, func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password), nil
}
