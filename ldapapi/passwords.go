package ldapapi

import (
	"bytes"
	"fmt"
	"strings"

	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
)

func sshaHash(passwordBytes, salt []byte) []byte {
	saltedPassword := append(passwordBytes, salt...)

	hasher := sha1.New()
	hasher.Write(saltedPassword)
	hashed := hasher.Sum(nil)

	return append(hashed, salt...)
}

func HashPasswordSSHA(password string, saltLength int) (string, error) {
	if len(password) > 6 && "{SSHA}" == password[:6] {
		return password, nil
	}

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate random salt")
	}
	passwdBytes := []byte(password)
	hashedPasswd := sshaHash(passwdBytes, salt)
	encodedHash := base64.StdEncoding.EncodeToString(hashedPasswd)
	return fmt.Sprintf("{SSHA}%s", encodedHash), nil
}

func VerifyHashSSHA(password, storedHash string) (bool, error) {
	if !bytes.HasPrefix([]byte(storedHash), []byte("{SSHA}")) {
		return false, fmt.Errorf("stored hash is not in SSHA format")
	}
	encodedHash := strings.TrimPrefix(storedHash, "{SSHA}")

	decoded, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode base64 hash: %w", err)
	}

	hashLen := sha1.Size
	if len(decoded) < hashLen {
		return false, fmt.Errorf("decoded hash is too short")
	}

	salt := decoded[hashLen:]
	hashBytes := decoded[:hashLen]

	pwdBytes := []byte(password)
	calculatedHash := sshaHash(pwdBytes, salt)[:hashLen]
	return bytes.Equal(calculatedHash, hashBytes), nil
}
