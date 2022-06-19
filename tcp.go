package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const (
	CONN_HOST = ""
	CONN_TYPE = "tcp"
)

func listen(c chan []byte, port string) {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn, c)
	}
}

func sendToSocket(data []byte, port string) {
	con, err := net.Dial("tcp", ":"+port)
	checkError(err)

	_, err = con.Write(data)
	checkError(err)

	res, err := ioutil.ReadAll(con)
	checkError(err)

	fmt.Println(string(res))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, c chan []byte) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Builds the message.

	c <- buf
	// Close the connection when you're done with it.
	conn.Close()
}
