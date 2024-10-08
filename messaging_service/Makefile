build: clean format_project
	mkdir -p bin
	go build -o bin/main main.go

run: build
	clear && \
	bin/main

clean:
	rm bin -rf

format_project:
	gofmt -s -w .
	go mod tidy
