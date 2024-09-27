package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gateway_service/internal"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyTok(t string) (*jwt.Token, error) {

	publicKey, err := getPublicKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error verifying token: %w", err)
	}

	return token, nil
}

func GenTok(name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": name,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	privateKey, err := getPrivateKey()
	if err != nil {
		return "", err
	}

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	keyData, ok := internal.GetCredentials()["pkey"]
	if !ok {
		return nil, fmt.Errorf("failed to read private key file")
	}

	keyString, ok := keyData.(string)
	if !ok {
		return nil, fmt.Errorf("invalid private key data")
	}

	block, _ := pem.Decode([]byte(keyString))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatal(block.Type)
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privKey1, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privKey1, nil
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}

	return rsaKey, nil
}

func getPublicKey() (*rsa.PublicKey, error) {
	keyData, ok := internal.GetCredentials()["pubkey"]
	if !ok {
		return nil, fmt.Errorf("failed to read public key")
	}

	keyString, ok := keyData.(string)
	if !ok {
		return nil, fmt.Errorf("invalid public key data")
	}

	block, _ := pem.Decode([]byte(keyString))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return publicKey, nil
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not a RSA public key")
	}

	return rsaPublicKey, nil
}
