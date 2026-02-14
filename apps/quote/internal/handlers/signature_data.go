package handlers

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strings"
)

func decodeSignatureData(signatureData string) ([]byte, error) {
	payload := strings.TrimSpace(signatureData)
	payload = strings.TrimPrefix(payload, "data:image/png;base64,")
	return base64.StdEncoding.DecodeString(payload)
}

func hasSignatureStroke(data []byte) bool {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return false
	}
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				return true
			}
		}
	}
	return false
}
