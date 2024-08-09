package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct {
}

func (fs *FileServer) start() {
	fmt.Println("FileServer started")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}

		// panic("Received more than 1000 bytes")
		fmt.Println(buf.Bytes())
		fmt.Printf("Received %d bytes over the network\n", n)
	}
}

func sendFile(size int) error {
	file := make([]byte, size)

	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))

	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))

	if err != nil {
		return err
	}

	fmt.Printf("Writtend %d bytes\n over the network \n", n)

	return nil

}

// main function
func main() {

	go func() {
		time.Sleep(4 * time.Second)
		sendFile(400000)
	}()

	server := &FileServer{}
	server.start()
}
