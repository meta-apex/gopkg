package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

func RsaEncrypt(origData, publicKey []byte) ([]byte, error) {
	pub, err := x509.ParsePKCS1PublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RsaDecrypt(ciphertext, privateKey []byte) ([]byte, error) {
	priv, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// RsaEncryptOAEP takes a message as a byte slice and an RSA public key, and returns the encrypted message as a byte slice.
func RsaEncryptOAEP(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	// Encrypt the message using RSA OAEP with SHA-256
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, msg, nil)
	if err != nil {
		// If there is an error encrypting the message, print the error message and return an empty byte slice
		return nil, err
	}

	// Return the encrypted message
	return ciphertext, nil
}

// RsaDecryptOAEP takes a message as a byte slice and an RSA private key, and returns the decrypted message as a byte slice.
func RsaDecryptOAEP(msg []byte, priv *rsa.PrivateKey) ([]byte, error) {
	// Decrypt the message using RSA OAEP with SHA-256
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, msg, nil)
	if err != nil {
		// If there is an error decrypting the message, print the error message and return an empty byte slice
		return nil, err
	}

	// Return the decrypted message
	return plaintext, nil
}

// PrivateKeyAsPEM takes a pointer to an RSA private key and returns its PEM encoding as a byte slice.
func PrivateKeyAsPEM(key *rsa.PrivateKey) []byte {
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	return privateKeyPEM
}

// PublicKeyAsPEM takes a pointer to an RSA private key and returns its corresponding public key's PEM encoding as a byte slice.
func PublicKeyAsPEM(key *rsa.PublicKey) []byte {
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(key),
		},
	)

	return publicKeyPEM
}

func PrivateBytesToPEM(buf []byte) []byte {
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: buf,
		},
	)

	return privateKeyPEM
}

func PublicBytesToPEM(buf []byte) []byte {
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: buf,
		},
	)

	return publicKeyPEM
}
