BINARY_NAME=pflags
GO_FILES=cli/pflags.go

build:
	go build -o ./out/$(BINARY_NAME) $(GO_FILES)

install:
	cp ./out/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	chmod +x /usr/local/bin/$(BINARY_NAME)

clean:
	go clean
	rm -f -r ./out

