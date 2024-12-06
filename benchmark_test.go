package gologger

import (
	"fmt"
	"testing"
)

func BenchmarkPrintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Printf("Processing item %d\n", i)
	}
}

func BenchmarkLoggerPrint(b *testing.B) {
	Start()
	defer Stop()

	for i := 0; i < b.N; i++ {
		Print(fmt.Sprintf("Processing item %d", i))
	}
}
