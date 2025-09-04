package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const privateKeyFileName = "private-key.pem"
const publicKeyFileName = "public-key.pem"

func main() {
	fmt.Println()
	privateKey, err := getPrivateKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting private key: %v\n", err)
		os.Exit(1)
	}

	jsonPubKey, err := getPublicKeyJson(privateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting public key JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(jsonPubKey)
	fmt.Println()

	signedJWT, err := generateSignedJWT(privateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating signed JWT: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(signedJWT)
	fmt.Println()
}

func generateSignedJWT(priv *rsa.PrivateKey) (string, error) {
	now := time.Now()
	issued := now.Unix()
	expires := now.AddDate(1, 0, 0).Unix()
	claims := jwt.MapClaims{
		"iat": issued,
		"exp": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(priv)
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	if _, err := os.Stat(privateKeyFileName); err == nil {
		return loadPrivateKey()
	} else if os.IsNotExist(err) {
		return createPrivateKey()
	} else {
		return nil, err
	}
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(privateKeyFileName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func createPrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	if err := writePrivateKey(privateKey); err != nil {
		return nil, err
	}

	if err := writePublicKey(privateKey); err != nil {
		return nil, err
	}

	return privateKey, nil
}

func writePrivateKey(priv *rsa.PrivateKey) error {
	keyFile, err := os.Create(privateKeyFileName)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}
	if err := pem.Encode(keyFile, pemBlock); err != nil {
		return err
	}

	return nil
}

func writePublicKey(priv *rsa.PrivateKey) error {
	publicKey := &priv.PublicKey

	keyFile, err := os.Create(publicKeyFileName)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	pemBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}
	if err := pem.Encode(keyFile, pemBlock); err != nil {
		return err
	}

	return nil
}

func getPublicKeyJson(priv *rsa.PrivateKey) (string, error) {
	publicKey := &priv.PublicKey
	keyData := x509.MarshalPKCS1PublicKey(publicKey)
	block := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyData,
	})

	keys := map[string]string{
		"key1": string(block),
	}

	jsonData, err := json.Marshal(keys)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
