BIN := tmp/tasks

.PHONY: build clean copyconf cov serve test

build:
	go build -o $(BIN)

clean:
	rm -f $(BIN)

copyconf:
	cp config/dev.json tmp/

cov:
	go test -coverprofile=tmp/coverage.out ./internal/...
	go tool cover -func=tmp/coverage.out

serve: build
	$(BIN) -path tmp/dev.json -task serve

test:
	go test ./internal/...