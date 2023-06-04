build:
	@go build -o .build/$(BINARY_NAME) -v
run:
	@go run main.go
clean:
	@rm -rf .build
	@rm -rf .cache
	@rm -rf .tmp
	@rm -rf .test
	@rm -rf .vendor
	@rm -rf .coverage
	@rm -rf .profile

test :
	@go test -v ./... -race