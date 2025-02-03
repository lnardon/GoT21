build:
	go build -o main *.go
run:
	make build && ./main
