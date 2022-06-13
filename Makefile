build:
	go build -trimpath -ldflags "-s -w" .

build-all: crypt_darwin_amd64 crypt_darwin_arm64 crypt_linux_amd64 crypt_windows_amd64 crypt_openbsd_amd64

crypt_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o $@ .
crypt_darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o $@ .
crypt_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o $@.
crypt_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o $@ .
crypt_openbsd_amd64:
	GOOS=openbsd GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o $@ .

clean:
	rm -f crypt crypt_*

.PHONY: build build-all clean
