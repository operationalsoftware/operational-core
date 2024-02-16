package migrate

import (
	"app/db"
	userModel "app/src/users/model"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func generateRandomPassword(length int) (string, error) {
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

func initialise() bool {
	db := db.UseDB()

	// start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	//
	// START OF INITIALISATION
	//

	fmt.Println("Initialising database...")

	// Create users table
	fmt.Print("Creating User table... ")
	_, err = db.Exec(`
CREATE TABLE User (
	UserID INTEGER PRIMARY KEY AUTOINCREMENT, 
	IsAPIUser BOOLEAN DEFAULT FALSE NOT NULL,
	Username TEXT NOT NULL UNIQUE,
	Email TEXT UNIQUE, 
	FirstName TEXT,
	LastName TEXT,
	Created DATETIME DEFAULT CURRENT_TIMESTAMP,
	LastLogin DATETIME DEFAULT NULL,
	HashedPassword TEXT NOT NULL,
	FailedLoginAttempts INTEGER DEFAULT 0 NOT NULL,
	LoginBlockedUntil DATETIME DEFAULT NULL,
	Roles JSON DEFAULT '[]' NOT NULL,
	UserData JSON DEFAULT '{}' NOT NULL
);
`)
	if err != nil {
		return false
	}
	fmt.Println("done")

	// add the system user with a random password
	fmt.Print("Creating system user... ")

	password, err := generateRandomPassword(24)
	if err != nil {
		panic(err)
	}

	var userToAdd = userModel.NewUser{
		Username:  "system",
		IsAPIUser: true,
		Password:  password,
		Roles:     []string{"User Admin"},
	}

	err = userModel.Add(tx, userToAdd)
	if err != nil {
		return false
	}

	fmt.Println("done")
	fmt.Println("System user details:\n\tusername: system\n\tpassword: " + password)

	//
	// END OF INITIALISATION
	//
	fmt.Print("Database initialised, committing changes...")
	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("done")

	return true
}
