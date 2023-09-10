.PHONY: test

test:
	go test -v ./...

test-coverage:
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	open ./coverage/coverage.html
