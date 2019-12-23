package security

import "encoding/base64"

func Base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}
func Base64Decode(src []byte) []byte {
	dst, _ := base64.StdEncoding.DecodeString(string(src))
	return dst
}
