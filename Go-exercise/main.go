package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"strings"
)

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers          = "0123456789"
	symbols          = "!@#$%^&*()-_=+,.?/:;{}[]~"
)

func generateRandomPassword(length int, includeNumbers, includeSymbols, includeUppercase bool) string {
	var chars string
	if includeNumbers {
		chars += numbers
	}
	if includeSymbols {
		chars += symbols
	}
	if includeUppercase {
		chars += uppercaseLetters
	}
	chars += lowercaseLetters

	var password strings.Builder
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		password.WriteByte(chars[randomIndex.Int64()])
	}
	return password.String()
}
func main() {
	lengthPtr := flag.Int("length", 12, "An integer specifying the length of the password")
	includeNumbersPtr := flag.Bool("includeNumbers", false, "Include numbers in the password")
	includeSymbolsPtr := flag.Bool("includeSymbols", false, "Include symbols in the password")
	includeUppercasePtr := flag.Bool("includeUppercase", false, "Include uppercase letters in the password")
	passwordTypePtr := flag.String("type", "random", "Type of password to generate: random, alphanumeric, or pin")

	flag.Parse()

	var passwordType string
	switch *passwordTypePtr {
	case "random":
		passwordType = "random"
	case "alphanumeric":
		passwordType = "alphanumeric"
	case "pin":
		passwordType = "pin"
	default:
		fmt.Println("Invalid password type. Please choose 'random', 'alphanumeric', or 'pin'.")
		return
	}

	var password string
	switch passwordType {
	case "random":
		password = generateRandomPassword(*lengthPtr, *includeNumbersPtr, *includeSymbolsPtr, *includeUppercasePtr)
	case "alphanumeric":
		password = generateRandomPassword(*lengthPtr, true, false, *includeUppercasePtr)
	case "pin":
		password = generateRandomPassword(6, true, false, false)
	}

	fmt.Println("Generated Password:", password)
}
