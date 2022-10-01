package util

import "encoding/base64"

func EncodeString(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
