BINARY_NAME=pflags
GO_FILES=cli/pflags/pflags.go

build:
	go build -o ./out/$(BINARY_NAME) $(GO_FILES)

install:
	cp ./out/$(BINARY_NAME) $(INSTALLATION_PATH)/$(BINARY_NAME)
	chmod +x $(INSTALLATION_PATH)/$(BINARY_NAME)

clean:
	go clean
	rm -f -r ./out

uninstall:
	rm -f $(INSTALLATION_PATH)/$(BINARY_NAME)