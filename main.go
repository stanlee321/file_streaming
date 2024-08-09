package main

import (
	"flag"
	"log"
	"os"

	"github.com/stanlee321/file_handler/filetransfer"
)

func main() {
	mode := flag.String("mode", "cli", "Mode of operation: cli or http")
	filePath := flag.String("file", "", "Path to the file to send")
	address := flag.String("address", "localhost:8080", "Address to send the file to")
	saveDir := flag.String("save-dir", os.TempDir(), "Directory to save received files")

	flag.Parse()

	switch *mode {
	case "cli":
		if *filePath != "" {
			// Sending mode
			if err := filetransfer.SendFile(*filePath, *address); err != nil {
				log.Fatalf("Failed to send file: %v", err)
			}
		} else {
			// Receiving mode
			filetransfer.StartServer(*saveDir)
		}
	case "http":
		// Start HTTP server with Gin
		filetransfer.StartHTTPServer(*saveDir)
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}
