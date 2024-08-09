# File Handler

This is a simple Go project that demonstrates how to send and receive files using Go. The project has two modes: CLI and HTTP.

## Getting Started

### Prerequisites

- Go (version 1.16 or higher)

### Installation

1. Initialize the Go module:

   ```sh
   go mod init file_handler
   ```

2. Download the dependencies:

   ```sh
   go mod tidy
   ```

### Running the Project


#### 2. **CLI Mode**:

- Flags are used to choose between sending and receiving (`mode` flag).
- If the `file` flag is provided, the program sends the specified file.
- If no `file` is specified, the server starts and listens for incoming files.

#### 3. **HTTP Mode**:

- The Gin framework is used to handle POST requests at the `/upload` endpoint.
- Files are uploaded using multipart form data and saved to the specified directory.

### How to Run:

#### CLI Mode:

- **Start the server to receive files:**
  ```sh
  go run main.go --mode cli --save-dir /path/to/save
  ```
- **Send a file:**
  ```sh
  go run main.go --mode cli --file /path/to/demo.mp4 --address localhost:8080
  ```

#### HTTP Mode:

- **Start the HTTP server:**
  ```sh
  go run main.go --mode http --save-dir /path/to/save
  ```
- **Upload a file via HTTP POST**:
  ```bash
  curl -X POST -F "file=@/path/to/demo.mp4" http://0.0.0.0:8080/upload
  ```

### Building the Project

To build the project, use the following command:

```sh
go build -o file_handler main.go
```
