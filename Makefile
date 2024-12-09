# Default value for count
COUNT ?= 1

benchmark:
	@go test -bench="Benchmark(Concurrent)" -count=$(COUNT) -benchmem