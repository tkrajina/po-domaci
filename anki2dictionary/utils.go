package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"unicode"
)

func firstLetter(str string) (rune, bool) {
	for _, r := range str {
		if unicode.IsLetter(r) {
			return r, true
		}
	}
	return ' ', false
}

func ignoreNonLetters(str string) string {
	var buf bytes.Buffer
	for _, r := range str {
		if unicode.IsLetter(r) {
			buf.WriteRune(r)
		} else {
			buf.WriteRune(' ')
		}
	}
	return strings.TrimSpace(buf.String())
}

func jsonizePrettified(i interface{}) string {
	bytes, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
