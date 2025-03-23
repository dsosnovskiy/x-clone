lint:
	golangci-lint version && \
	golangci-lint run ./... --out-format colored-line-number