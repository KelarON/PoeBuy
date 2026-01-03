package utils

import (
	"encoding/base64"
	"os/user"
)

// Encrypt encrypts a string using the user's UID as the key. The encrypted string is then encoded in Base64.
func Encrypt(input string) (string, error) {

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	key := usr.Uid

	output := make([]byte, len(input))
	keyLen := len(key)

	for i := range input {
		output[i] = input[i] ^ key[i%keyLen]
	}

	// Encode the encrypted byte slice to Base64
	encoded := base64.StdEncoding.EncodeToString(output)
	return encoded, nil
}

// Decrypt decrypts a Base64 encoded string using the user's UID as the key.
func Decrypt(encoded string) (string, error) {

	// Decode the Base64 encoded string
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	key := usr.Uid

	output := make([]byte, len(decoded))
	keyLen := len(key)

	for i := range decoded {
		output[i] = decoded[i] ^ key[i%keyLen]
	}

	return string(output), nil
}
