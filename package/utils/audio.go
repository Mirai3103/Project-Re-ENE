package utils

import (
	"encoding/base64"
	"io"
)

func ReaderToBase64(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
