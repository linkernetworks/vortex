package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// SHA256String will do SHA256 String
func SHA256String(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	md := hash.Sum(nil)
	return fmt.Sprintf("%s", hex.EncodeToString(md))
}
