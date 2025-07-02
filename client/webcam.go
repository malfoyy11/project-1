package webcam

import (
    "gocv.io/x/gocv"
)

// CaptureSnapshot captures a webcam image and saves it as a JPG
func CaptureSnapshot(filename string) error {
    webcam, err := gocv.OpenVideoCapture(0)
    if err != nil {
        return err
    }
    defer webcam.Close()

    img := gocv.NewMat()
    defer img.Close()

    if ok := webcam.Read(&img); !ok {
        return err
    }

    if img.Empty() {
        return err
    }

    return gocv.IMWrite(filename, img)
}
