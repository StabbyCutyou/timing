.PHONY: test bench build

build:
	go build

test:
	go test

bench:
	go test -bench=. -benchmem