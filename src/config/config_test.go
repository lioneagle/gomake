package config

import (
	"bytes"
	"testing"
)

func TestBytesEqual(t *testing.T) {
	if !bytes.Equal([]byte("abc"), []byte("abc")) {
		t.Errorf("bytes.Equal failed")
	}
}

func BenchmarkBytesEqual(b *testing.B) {
	b.StopTimer()
	s1 := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s2 := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		bytes.Equal(s1, s2)
	}
}

func BenchmarkBytesEqualFold(b *testing.B) {
	b.StopTimer()
	s1 := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s2 := []byte("aBCDEFGHIJKLMNOPQRSTUVWXYz")

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		bytes.EqualFold(s1, s2)
	}
}
