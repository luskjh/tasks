BIN := tmp/tasks

.PHONY: build clean copyconf serve

build:
	go build -o $(BIN)

clean:
	rm -f $(BIN)

copyconf:
	cp config/dev.json tmp/

serve: build
	$(BIN) -path tmp/dev.json -task serve