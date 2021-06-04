name = bin/idea-activation

all: macos windows macos-arm linux

macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${name}-macos ./...

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${name}-windows.exe ./...

macos-arm:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${name}-macos-arm64 ./...

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${name}-linux ./...
