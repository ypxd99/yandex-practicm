package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"sync"

	"github.com/pkg/errors"
)

var (
	onceRSA     sync.Once
	rsaKey      *RSA
	publicPath  = "public.pem"
	privatePath = "private.pem"
)

type RSAEncryptor interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type RSA struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (r *RSA) Encrypt(data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, r.publicKey, data, nil)
}

func (r *RSA) Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, r.privateKey, data, nil)
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.WithMessage(err, "error occurred while reading public key file")
	}
	block, _ := pem.Decode(file)
	if block == nil {
		return nil, errors.WithMessage(errors.New("cant decode PEM"), "error occured while decoding public PEM")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.WithMessage(err, "error occurred while parsing public key")
	}
	return pub.(*rsa.PublicKey), nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.WithMessage(err, "error occurred while reading private key file")
	}
	block, _ := pem.Decode(file)
	if block == nil {
		return nil, errors.WithMessage(errors.New("cant decode PEM"), "error occured while decoding private PEM")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.WithMessage(err, "error occurred while parsing private key")
	}
	return priv, nil
}

func GetRSA() *RSA {
	onceRSA.Do(func() {
		privKey, err := loadPrivateKey(privatePath)
		if err != nil {
			log.Fatalf("error occurred while load private key: %s", err)
		}

		pubKey, err := loadPublicKey(publicPath)
		if err != nil {
			log.Fatalf("error occurred while load public key: %s", err)
		}

		rsaKey = &RSA{
			privateKey: privKey,
			publicKey:  pubKey,
		}
	})

	if rsaKey == nil {
		log.Fatal("nil rsa")
	}

	return rsaKey
}

func GenerateRSA() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		GetLogger().Errorf("error occurred while generating rsa key: %s", err)
		return
	}

	privateFile, _ := os.Create("rsa_private.pem")
	defer privateFile.Close()
	pem.Encode(privateFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKey := &privateKey.PublicKey

	pubASN1, _ := x509.MarshalPKIXPublicKey(publicKey)
	publicFile, _ := os.Create("rsa_public.pem")
	defer publicFile.Close()
	pem.Encode(publicFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
}
