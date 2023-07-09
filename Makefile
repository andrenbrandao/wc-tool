build:
	go build -o bin/ccwc cmd/ccwc/main.go

clean:
	rm -rf bin/

test:
	cd cmd/ccwc/ && pwd && go test -v