package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"log"
)

func main() {

	hash := func(b []byte) []byte {
		h := sha256.New()
		h.Write(b)
		return h.Sum(nil)
	}

	secretKey := newSecret(32)
	secretKey = hash(secretKey)
	plaintext := []byte("Hello, world.")

	// ciphertext is prepended with a 12-byte nonce
	ciphertext, err := encryptValue(secretKey, plaintext)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("generated ciphertext with len %v\n", len(ciphertext))

	recovered, err := decryptValue(secretKey, ciphertext)
	if err != nil {
		log.Fatal(err)
	}

	if bytes.Compare(recovered, plaintext) != 0 {
		log.Fatal("recovered text does not equal plaintext")
	}

}

const (
	nonceLen = 12
)

func encryptValue(enc []byte, v []byte) ([]byte, error) {
	block, err := aes.NewCipher(enc)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, nonceLen)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, v, nil)
	return append(nonce, ciphertext...), nil
}

func decryptValue(enc []byte, v []byte) ([]byte, error) {
	block, err := aes.NewCipher(enc)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := v[0:nonceLen]
	ciphertext := v[nonceLen:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func newSecret(size int) []byte {
	k := make([]byte, size)
	rand.Read(k)
	return k
}
