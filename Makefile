build:
	@go build -o bin/receipt-processor-challenge

run: build
	@./bin/receipt-processor-challenge