# Default value for count
COUNT ?= 1

benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem -count=$(COUNT)