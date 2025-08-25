build:
	@go build -o bin/paysplit

run: build
	@./bin/paysplit

test:
	@go test ./...s


