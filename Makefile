.PHONY: build-sim run-sim

build-sim:
	go build -ldflags "-s -w" -o ./build/orb-sim ./cmd

run-sim: build-sim
	./build/orb-sim
