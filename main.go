package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

type FileServer struct{}

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
	defer conn.Close()

	var size int64
	err := binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new file to save the received data
	tempDir := "./tmp" // " os.TempDir()"
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	receivedFilePath := filepath.Join(tempDir, "received_video.mp4")
	file, err := os.Create(receivedFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Copy the data from the connection to the file
	n, err := io.CopyN(file, conn, size)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received %d bytes and saved to %s\n", n, receivedFilePath)
}

func sendFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	size := fileInfo.Size()

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the file size first
	err = binary.Write(conn, binary.LittleEndian, size)
	if err != nil {
		return err
	}

	// Send the file data
	n, err := io.CopyN(conn, file, size)
	if err != nil {
		return err
	}

	fmt.Printf("Sent %d bytes over the network from file %s\n", n, filePath)
	return nil
}

func main() {
	fileName := "XVR_ch1_main_20210910141900_20210910142500.mp4"
	// Simulate sending a file after a delay
	go func() {
		time.Sleep(4 * time.Second)
		err := sendFile(fileName) // Change this to the actual path of your MP4 file
		if err != nil {
			log.Fatal(err)
		}
	}()

	server := &FileServer{}
	server.start()
}
