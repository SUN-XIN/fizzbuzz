package main

import (
	"testing"
)

var (
	size = 1000
	cr   clientRequest
)

func init() {
	cr = clientRequest{
		String1: "fizz",
		String2: "buzz",
		Int1:    3,
		Int2:    5,
		Limit:   size,
	}
}

func BenchmarkProcess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		processRequest(&cr)
	}
}

func BenchmarkProcessBis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		processRequestBis(&cr)
	}
}
