export GO111MODULE=on
install:
	@go mod tidy
	# @go build -v -o $(NAME) $(MAIN)
