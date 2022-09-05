package imapdump

import "strings"

var (
	sanitizeMessageIDChars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@.-_")
)

func SanitizeMessageID(s string) string {
	s = strings.TrimSuffix(strings.TrimPrefix(s, "<"), ">")
	u := []byte(s)
outer:
	for i, c := range u {
		for _, c1 := range sanitizeMessageIDChars {
			if c1 == c {
				continue outer
			}
		}
		u[i] = '-'
	}
	return string(u)
}
