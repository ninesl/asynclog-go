# Default value for count
COUNT ?= 1

benchmark:
	@go test -bench="Benchmark(ConcurrentGo|Fmt)" -count=$(COUNT) -benchmem