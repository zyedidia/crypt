build:
	go build -trimpath -ldflags "-s -w" .

.PHONY: build
