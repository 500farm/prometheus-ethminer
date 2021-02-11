.PHONY: build clean

bin/ethminer_exporter:
	@GOOS=linux GOARCH=amd64 go build -o bin/ethminer_exporter src/*.go

build: bin/ethminer_exporter

clean:
	@rm -rf ./bin
