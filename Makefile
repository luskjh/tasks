BIN := tmp/tasks

.PHONY: build clean run

build:
	go build -o $(BIN)

clean:
	rm -f $(BIN)

run: build
	$(BIN)