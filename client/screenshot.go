package screenshot

import (
    "image/png"
    "os"

    "github.com/kbinani/screenshot"
)

func CaptureAndSave(filename string) error {
    bounds := screenshot.GetDisplayBounds(0)
    img, err := screenshot.CaptureRect(bounds)
    if err != nil {
        return err
    }

    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    return png.Encode(file, img)
}
