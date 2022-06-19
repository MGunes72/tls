package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Message struct {
	CipherText string `json:"cipherText"`
	Mac        string `json:"mac"`
}

func createMac(msg string, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	msgBytes, _ := hex.DecodeString(msg)
	mac.Write(msgBytes)

	var message Message
	message.CipherText = msg
	message.Mac = hex.EncodeToString(mac.Sum(nil))

	fmt.Println("Sended CipherText: ", msg)
	fmt.Println("Mac:", message.Mac)

	return encodeMessageToBytes(message)
}

func verifyMac(msgAndMac, key []byte) (string, bool, error) {
	msg := decodeMessageBytes(msgAndMac)
	sig, err := hex.DecodeString(msg.Mac)
	if err != nil {
		return "", false, err
	}

	mac := hmac.New(sha256.New, key)
	msgBytes, _ := hex.DecodeString(msg.CipherText)
	mac.Write(msgBytes)

	return msg.CipherText, hmac.Equal(sig, mac.Sum(nil)), nil
}

func encodeMessageToBytes(p Message) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func decodeMessageBytes(b []byte) Message {
	tmpbuff := bytes.NewBuffer(b)
	tmpstruct := new(Message)

	gobobjdec := gob.NewDecoder(tmpbuff)
	gobobjdec.Decode(tmpstruct)

	return *tmpstruct
}
