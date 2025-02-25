NAME := latai
VERSION := 0.1.0
DIST_DIR := dist

.PHONY: build


build: build-mac build-linux build-win


build-mac:
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(NAME)-$(VERSION)-mac-amd64
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(NAME)-$(VERSION)-mac-arm64


build-linux:
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(NAME)-$(VERSION)-linux-amd64
	GOOS=linux GOARCH=arm64 go build -o $(DIST_DIR)/$(NAME)-$(VERSION)-linux-arm64


build-win:
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(NAME)-$(VERSION)-win-amd64.exe


clean:
	rm -rf $(DIST_DIR)