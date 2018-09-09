package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
)

func ComputeSHA256Checksum(reader io.Reader) string {
	h := sha256.New()
	if _, err := io.Copy(h, reader); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
