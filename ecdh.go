package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Hello struct {
	DhPub *ecdsa.PublicKey `json:"dhPub"`
	Id    []byte           `json:"Id"`
	DsPub *ecdsa.PublicKey `json:"dsPub"`
	Sig   []byte           `json:"sig"`
}

func clientHello(cx chan []byte, cx2 chan []byte, id string) {
	gob.Register(Hello{})
	gob.Register(elliptic.P256())
	priv, pub := generateKeyECDSA()
	privKeyDS, publicKeyECDSA := generateKeyECDSA()

	var cHello Hello
	cHello.Id = []byte(id)

	pubBytes, _ := json.Marshal(pub)
	pubBytes = append(pubBytes, cHello.Id...)

	cHello.DhPub = pub
	cHello.DsPub = publicKeyECDSA

	sig, _, err := signMessage(pubBytes, privKeyDS)
	if err != nil {
		fmt.Println("err", err)
	}

	cHello.Sig = sig

	msg := encodeToBytes(cHello)

	fmt.Println("Sended message to B ", hex.EncodeToString(msg))

	sendToSocket(msg, "3333")

	c := make(chan []byte)
	go listen(c, "3334")

	receivedMessage := <-c

	sHello := decodeBytes(receivedMessage)

	pubSBytes, _ := json.Marshal(sHello.DhPub)
	pubSBytes = append(pubSBytes, sHello.Id...)

	hashB := sha256.Sum256(pubSBytes)

	fmt.Println("Received message by B: ", hex.EncodeToString(receivedMessage))

	err = verifySignature(sHello.Sig, hashB, sHello.DsPub)
	if err != nil {
		fmt.Println("err", err)
	}

	secret := computeSecret(sHello.DhPub, priv)
	fmt.Println("secret shared: ", hex.EncodeToString(secret))
	cx <- secret
	cx2 <- sHello.Id
}

func serverHello(cx chan []byte, cx2 chan []byte, id string) {
	gob.Register(Hello{})
	gob.Register(elliptic.P256())
	c := make(chan []byte)
	go listen(c, "3333")

	time.Sleep(10 * time.Second)

	receivedMessage := <-c

	cHello := decodeBytes(receivedMessage)

	pubCBytes, _ := json.Marshal(cHello.DhPub)
	pubCBytes = append(pubCBytes, cHello.Id...)

	hashA := sha256.Sum256(pubCBytes)

	fmt.Println("Received message by A: ", hex.EncodeToString(receivedMessage))

	err := verifySignature(cHello.Sig, hashA, cHello.DsPub)
	if err != nil {
		fmt.Println("err", err)
	}

	priv, pub := generateKeyECDSA()
	privKeyDS, publicKeyECDSA := generateKeyECDSA()

	var sHello Hello
	sHello.Id = []byte(id)

	pubBytes, _ := json.Marshal(pub)
	pubBytes = append(pubBytes, sHello.Id...)

	sHello.DhPub = pub
	sHello.DsPub = publicKeyECDSA

	sig, _, err := signMessage(pubBytes, privKeyDS)
	if err != nil {
		fmt.Println("err", err)
	}

	sHello.Sig = sig
	msg := encodeToBytes(sHello)

	fmt.Println("Sended message to A ", hex.EncodeToString(msg))

	sendToSocket(msg, "3334")

	secret := computeSecret(cHello.DhPub, priv)
	fmt.Println("secret shared: ", hex.EncodeToString(secret))
	cx <- secret
	cx2 <- cHello.Id
}

func computeSecret(pubA *ecdsa.PublicKey, priv *ecdsa.PrivateKey) []byte {
	a, _ := pubA.Curve.ScalarMult(pubA.X, pubA.Y, priv.D.Bytes())

	sharedSecret := sha256.Sum256(a.Bytes())

	return sharedSecret[:]
}

func encodeToBytes(p Hello) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func decodeBytes(b []byte) Hello {
	tmpbuff := bytes.NewBuffer(b)
	tmpstruct := new(Hello)

	gobobjdec := gob.NewDecoder(tmpbuff)
	gobobjdec.Decode(tmpstruct)

	return *tmpstruct
}
