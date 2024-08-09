package filetransfer

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// ReceiveFile handles receiving a file over TCP and saving it to the specified path.
func ReceiveFile(conn net.Conn, savePath string) error {
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

// SendFile handles sending a file over TCP to the specified address.
func SendFile(filePath, address string) error {
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
func StartServer(saveDir string) {
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
			if err := ReceiveFile(conn, savePath); err != nil {
				log.Printf("Error receiving file: %v", err)
			}
		}()
	}
}

// HTTP Mode: Gin server for handling file uploads
func StartHTTPServer(saveDir string) {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(400, "Bad request: %v", err)
			return
		}
		// If saveDir does not exists, create it
		if _, err := os.Stat(saveDir); os.IsNotExist(err) {
			os.MkdirAll(saveDir, os.ModePerm)
		}

		// Set current date as YYYYMMDD
		currentDate := time.Now().Format("20060102")
		extenstion := filepath.Ext(file.Filename)
		// Remove extension from filename
		fileNameWithoutExt := file.Filename[0 : len(file.Filename)-len(extenstion)]

		// Save file with current date and received_video suffix
		fileName := fmt.Sprintf("%s_%s_%s%s", currentDate, fileNameWithoutExt, "received_video", extenstion)
		savePath := filepath.Join(saveDir, fileName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.String(500, "Failed to save file: %v", err)
			return
		}

		c.String(200, "File uploaded successfully")
	})

	r.Run(":8080")
}
