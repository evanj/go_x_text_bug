package main

import (
	"strconv"
	"testing"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

func TestNFCIterBug(t *testing.T) {
	const maxIterations = 20
	const badInput = "\xf0\xd9\x95"
	nfcString := norm.NFC.String(badInput)

	for i, b := range []byte(badInput) {
		t.Logf("badInput byte i=%d b=0x%x", i, b)
	}

	for i, r := range nfcString {
		t.Logf("nfcString rune i=%d r=0x%x", i, r)
	}

	iter := norm.Iter{}
	iter.InitString(norm.NFC, nfcString)
	i := 0
	for !iter.Done() {
		bytes := iter.Next()
		t.Logf("norm Iter i=%d bytes=%#v", i, bytes)
		i += 1
		if i > maxIterations {
			t.Fatalf("stopping after %d iterations to avoid infinite loop", maxIterations)
		}
	}
}

// Compare norm.NFC.Bytes to norm.Iter.
func FuzzNFCIterator(f *testing.F) {
	f.Add("")
	f.Add("ascii")
	f.Add("e\u0301 decomposed")

	f.Fuzz(func(t *testing.T, s string) {
		// check UTF-8 valid strings only: fuzzing appears to work
		// if !utf8.ValidString(s) {
		// 	return
		// }
		normalized := string(norm.NFC.String(s))
		runes := []rune(normalized)

		iter := norm.Iter{}
		iter.InitString(norm.NFC, normalized)
		runeI := 0
		for !iter.Done() {
			runeBytes := iter.Next()
			if len(runeBytes) == 0 {
				t.Fatalf("iter.Next() returned empty byte slice for s=%#v %s",
					s, strconv.QuoteToASCII(s))
			}
			for len(runeBytes) > 0 {
				rune, runeLen := utf8.DecodeRune(runeBytes)
				runeBytes = runeBytes[runeLen:]
				if runes[runeI] != rune {
					t.Fatalf("s=%#v %s: runes[runeI=%d]=0x%x iter returned 0x%x",
						s, strconv.QuoteToASCII(s), runeI, runes[runeI], rune)
				}
				runeI++
			}
		}
		if runeI != len(runes) {
			t.Fatalf("s=%#v %s: expected %d runes, got %d",
				s, strconv.QuoteToASCII(s), len(runes), runeI)
		}
	})
}
