PROJECT:=ferry

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ferry main.go && upx -9 ferry
build-sqlite:
	go build -tags sqlite3 -o ferry main.go
#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t ferry:latest
