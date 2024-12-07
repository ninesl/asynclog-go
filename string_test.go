package gologger_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// resembles the debugInfo struct from logger.go
type DebugInfo struct {
	file string
	line int
}

func (info *DebugInfo) String() string {
	return fmt.Sprintf("%s:%d", info.file, info.line)
}

func generateInput(size int) []string {
	input := make([]string, size)
	for i := 0; i < size; i++ {
		input[i] = "test"
	}
	return input
}

func BenchmarkStringConcat(b *testing.B) {
	info := &DebugInfo{file: "test.go", line: 42}
	msg := "test message"

	b.Run("Sprintf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%s %s", info, msg)
		}
	})

	b.Run("+", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = info.String() + " " + msg
		}
	})

	b.Run("StringBuilder", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var sb strings.Builder
			sb.WriteString(info.String())
			sb.WriteString(" ")
			sb.WriteString(msg)
			_ = sb.String()
		}
	})

	b.Run("StringBuilderPrealloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var sb strings.Builder
			sb.Grow(len(info.String()) + 1 + len(msg))
			sb.WriteString(info.String())
			sb.WriteString(" ")
			sb.WriteString(msg)
			_ = sb.String()
		}
	})
}

func BenchmarkPrintMethods(b *testing.B) {
	sizes := []int{1, 2, 3, 4, 5, 10, 20, 50, 100, 200}

	for _, size := range sizes {
		input := generateInput(size)

		b.Run("+="+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var m string
				for _, arg := range input {
					m += arg
				}
			}
		})

		b.Run("join"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = strings.Join(input, "")
			}
		})

		b.Run("growbuilder"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				totalLen := 0
				for _, s := range input {
					totalLen += len(s)
				}
				var sb strings.Builder
				sb.Grow(totalLen)
				for _, arg := range input {
					sb.WriteString(arg)
				}
				_ = sb.String()
			}
		})

		b.Run("builder"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				totalLen := 0
				for _, s := range input {
					totalLen += len(s)
				}
				var sb strings.Builder
				for _, arg := range input {
					sb.WriteString(arg)
				}
				_ = sb.String()
			}
		})
	}
}
