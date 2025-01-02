package keygenerator

import "golang.org/x/crypto/argon2"

func GenerateKey(generator []byte, key []byte) []byte {
	return argon2.IDKey(generator, key, 1, 64*1024, 4, 32)
}
