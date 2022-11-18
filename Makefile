COVERPROFILE = ./coverage.out

run:
	go run cmd/brainfuck/main.go -i example.bf

test:
	go test -v --race ./... -coverprofile=$(COVERPROFILE)

test-coverage: test
	go tool cover -html=$(COVERPROFILE)
