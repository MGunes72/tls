package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

func main() {
	mode := os.Args[1]
	id := os.Args[2]
	pl := os.Args[3]

	if mode == "A" {
		c := make(chan []byte)
		c2 := make(chan []byte)
		go clientHello(c, c2, id)

		secretKey := <-c
		idServer := <-c2

		keyReceive := createDerivedKey(secretKey, idServer)
		keySend := createDerivedKey(secretKey, []byte(id))

		idMac := append([]byte(id), idServer...)
		keyMac := createDerivedKey(secretKey, idMac)

		fmt.Println("Key receive A: ", hex.EncodeToString(keyReceive))
		fmt.Println("Key send A: ", hex.EncodeToString(keySend))
		fmt.Println("Key mac A: ", hex.EncodeToString(keyMac))

		ciphertext, _ := encrypt(keySend[:], pl)
		time.Sleep(10 * time.Second)

		fmt.Println("Sended message: ", pl)

		msg := createMac(ciphertext, keyMac)
		sendToSocket(msg, "3335")

		cx := make(chan []byte)
		go listen(cx, "3336")

		cipherReceive := <-cx
		cipherString, verified, err := verifyMac(cipherReceive, keyMac)
		if err != nil {
			panic(err)
		}

		if !verified {
			fmt.Println("Mac not verified")
		}

		plaintextReceive, err := decrypt(keyReceive[:], cipherString)
		if err != nil {
			panic(err)
		}
		fmt.Println("Received ciphertext: ", cipherString)
		fmt.Println("Plaintext:", plaintextReceive)

	} else if mode == "B" {
		c := make(chan []byte)
		c2 := make(chan []byte)
		go serverHello(c, c2, id)

		secretKey := <-c
		idClient := <-c2

		keyReceive := createDerivedKey(secretKey, idClient)
		keySend := createDerivedKey(secretKey, []byte("Bob"))

		idMac := append(idClient, []byte("Bob")...)
		keyMac := createDerivedKey(secretKey, idMac)

		fmt.Println("Key receive B:", hex.EncodeToString(keyReceive))
		fmt.Println("Key send B:", hex.EncodeToString(keySend))
		fmt.Println("Key mac B:", hex.EncodeToString(keyMac))

		cx := make(chan []byte)
		go listen(cx, "3335")

		ciphertext := <-cx
		cipherString, verified, err := verifyMac(ciphertext, keyMac)
		if err != nil {
			panic(err)
		}

		if !verified {
			fmt.Println("Mac not verified")
		}

		plaintext, err := decrypt(keyReceive[:], cipherString)
		if err != nil {
			fmt.Println("Err: ", err)
			panic(err)
		}
		fmt.Println("Received ciphertext: ", cipherString)
		fmt.Println("Plaintext:", plaintext)

		fmt.Println("Sended message: ", pl)
		ciphertextSend, _ := encrypt(keySend[:], pl)
		time.Sleep(10 * time.Second)

		msg := createMac(ciphertextSend, keyMac)
		sendToSocket(msg, "3336")
	}

}
