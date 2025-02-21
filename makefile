build:
	@go mod tidy
	@go build -v -o ./bin/wild . 

run: build
	@./bin/wild	
