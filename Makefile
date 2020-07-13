PROJECT:=ferry

.PHONY: build
build:
	CGO_ENABLED=0 go build -o ferry main.go
build-sqlite:
	go build -tags sqlite3 -o ferry main.go
#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t ferry:latest
