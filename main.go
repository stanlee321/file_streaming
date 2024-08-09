package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// Core function to receive a file over TCP
func receiveFile(conn net.Conn, savePath string) error {
	defer conn.Close()

	var size int64
	err := binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		return err
	}

	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := io.CopyN(file, conn, size)
	if err != nil {
		return err
	}

	fmt.Printf("Received %d bytes and saved to %s\n", n, savePath)
	return nil
}

// Core function to send a file over TCP
func sendFile(filePath, address string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	size := fileInfo.Size()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = binary.Write(conn, binary.LittleEndian, size)
	if err != nil {
		return err
	}

	n, err := io.CopyN(conn, file, size)
	if err != nil {
		return err
	}

	fmt.Printf("Sent %d bytes from file %s to %s\n", n, filePath, address)
	return nil
}

// CLI Mode: FileServer that listens for incoming connections
func startServer(saveDir string) {
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

		go func() {
			savePath := filepath.Join(saveDir, "received_video.mp4")
			if err := receiveFile(conn, savePath); err != nil {
				log.Printf("Error receiving file: %v", err)
			}
		}()
	}
}

// HTTP Mode: Gin server for handling file uploads
func startHTTPServer(saveDir string) {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(400, "Bad request: %v", err)
			return
		}
		// Set current date as YYYYMMDD
		currentDate := time.Now().Format("20060102")
		fileName := fmt.Sprintf("%s_%s_%s", currentDate, file.Filename, "received_video.mp4")
		savePath := filepath.Join(saveDir, fileName)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.String(500, "Failed to save file: %v", err)
			return
		}

		c.String(200, "File uploaded successfully")

	})

	r.Run(":8080")
}

func main() {
	// CLI flags
	mode := flag.String("mode", "cli", "Mode of operation: cli or http")
	filePath := flag.String("file", "", "Path to the file to send")
	address := flag.String("address", "localhost:8080", "Address to send the file to")
	saveDir := flag.String("save-dir", os.TempDir(), "Directory to save received files")

	flag.Parse()

	switch *mode {
	case "cli":
		if *filePath != "" {
			// Sending mode
			if err := sendFile(*filePath, *address); err != nil {
				log.Fatalf("Failed to send file: %v", err)
			}
		} else {
			// Receiving mode
			startServer(*saveDir)
		}
	case "http":
		// Start HTTP server with Gin
		startHTTPServer(*saveDir)
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}
