package util

import "encoding/base64"

func DecodeB64(value string) string {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "(decode error)"
	}
	return string(decoded)
}
