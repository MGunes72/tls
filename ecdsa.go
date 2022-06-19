package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
)

func signMessage(msg []byte, privateKey *ecdsa.PrivateKey) ([]byte, [32]byte, error) {
	hash := sha256.Sum256(msg)

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		return sig, hash, err
	}

	return sig, hash, nil
}

func verifySignature(sig []byte, hash [32]byte, publicKey *ecdsa.PublicKey) error {
	valid := ecdsa.VerifyASN1(publicKey, hash[:], sig)
	if !valid {
		return errors.New("failed to verify signature")
	}
	fmt.Println("signature verified:", valid)
	return nil
}

func generateKeyECDSA() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey
}
