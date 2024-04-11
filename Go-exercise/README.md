# go-exercises
 # CLI Password Generator in Go

## Overview
In this programming exercise, you are tasked with creating a Command Line Interface (CLI) application in Go that generates passwords based on user preferences. The application should allow users to specify the type and complexity of the password they need, offering various options such as length, inclusion of numbers, symbols, and uppercase letters.

## Goals
- Familiarize yourself with Go's standard library, particularly the packages related to string manipulation, random number generation, and CLI argument parsing.
- Practice building CLI applications that interact with users through the command line.
- Understand the principles behind generating random data securely in Go.

## Requirements
### CLI Arguments
The application should accept the following command-line arguments:
- `-length`: An integer specifying the length of the password (default to 12 if not provided).
- `-includeNumbers`: A boolean flag indicating whether to include numbers in the password.
- `-includeSymbols`: A boolean flag indicating whether to include symbols (e.g., !@#$%) in the password.
- `-includeUppercase`: A boolean flag indicating whether to include uppercase letters in the password.
- `-type`: A string indicating the type of password to generate. Valid options are random (default), alphanumeric, and pin.
  - `random`: Any characters can be included based on the flags provided.
  - `alphanumeric`: Only letters and numbers are included.
  - `pin`: Only numbers are included, and the default length is set to 6 unless specified.

### Secure Randomness
Use Go's `crypto/rand` package to generate secure random passwords.

### Error Handling
The application should gracefully handle any user input errors, such as invalid argument values.

### Help Command
Implement a help command or flag that displays usage instructions.

## Usage
```bash
# Generate a random password with default settings
./password_generator

# Generate a password of length 16 with numbers and symbols included
./password_generator -length=16 -includeNumbers -includeSymbols

# Generate an alphanumeric password with uppercase letters included
./password_generator -type=alphanumeric -includeUppercase
