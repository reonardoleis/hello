package utils

import "bytes"

func SanitizeBuffer(buf []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if r == '\x00' {
			return -1
		}

		return r
	}, buf)
}
