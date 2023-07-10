

format-tools:
	go install mvdan.cc/gofumpt@v0.3.1
	gofumpt -l -w .

lint: format-tools
	golangci-lint run

format: format-tools
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofumpt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/jackalLabs/jackalgo


.PHONY: lint format format-tools
