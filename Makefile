build: get-deps
	go build cmd/*.go

get-deps:
	go get github.com/gorilla/websocket

test:
	go test -v .
