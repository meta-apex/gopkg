package crypto

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"math/big"
)

// ECDH The main interface for ECDH key exchange.
type ECDH interface {
	// GenerateKey 生成公私钥对
	GenerateKey() error
	// MarshalPublicKey 将公钥转换为字节
	MarshalPublicKey() []byte
	// MarshalPrivateKey 将私钥转换为字节
	MarshalPrivateKey() []byte
	// GenerateSharedSecret 生成共享密钥
	GenerateSharedSecret([]byte) ([]byte, []byte, error)
}

type (
	curveECDHPublicKey struct {
		elliptic.Curve
		X, Y *big.Int
	}

	curveECDHPrivateKey struct {
		D []byte
	}

	curveECDH struct {
		curve      elliptic.Curve
		publicKey  curveECDHPublicKey
		peerPubKey curveECDHPublicKey
		privateKey curveECDHPrivateKey
	}
)

func NewCurveKey() ECDH {
	return &curveECDH{
		curve: elliptic.P256(),
	}
}

func (c *curveECDH) GenerateKey() error {
	d, x, y, err := elliptic.GenerateKey(c.curve, rand.Reader)
	if err != nil {
		return err
	}

	c.privateKey.D = d
	c.publicKey.X = x
	c.publicKey.Y = y
	c.publicKey.Curve = c.curve
	c.peerPubKey.Curve = c.curve

	return nil
}

func (c *curveECDH) MarshalPublicKey() []byte {
	return elliptic.Marshal(c.curve, c.publicKey.X, c.publicKey.Y)
}

func (c *curveECDH) MarshalPrivateKey() []byte {
	return c.privateKey.D
}

func (c *curveECDH) GenerateSharedSecret(buf []byte) ([]byte, []byte, error) {
	c.peerPubKey.X, c.peerPubKey.Y = elliptic.Unmarshal(c.curve, buf)
	x, _ := c.curve.ScalarMult(c.peerPubKey.X, c.peerPubKey.Y, c.privateKey.D)
	return x.Bytes(), x.Bytes(), nil
}

type rsaECDH struct {
	publicKey  *rsa.PublicKey
	peerPubKey *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewRsaKey() ECDH {
	return &rsaECDH{
		publicKey:  nil,
		peerPubKey: nil,
		privateKey: nil,
	}
}

func (c *rsaECDH) GenerateKey() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		// If there is an error generating the key, print the error message and return empty private and public keys
		fmt.Println("Error generating RSA key:", err)
		return err
	}

	c.privateKey = privateKey
	c.publicKey = &privateKey.PublicKey
	return nil
}

func (c *rsaECDH) MarshalPublicKey() []byte {
	return x509.MarshalPKCS1PublicKey(c.publicKey)
}

func (c *rsaECDH) MarshalPrivateKey() []byte {
	return x509.MarshalPKCS1PrivateKey(c.privateKey)
}

func (c *rsaECDH) GenerateSharedSecret(buf []byte) (pri []byte, pub []byte, err error) {
	c.peerPubKey, err = x509.ParsePKCS1PublicKey(buf)
	if err != nil {
		return
	}

	pri = x509.MarshalPKCS1PrivateKey(c.privateKey)
	pub = x509.MarshalPKCS1PublicKey(c.peerPubKey)
	return
}
