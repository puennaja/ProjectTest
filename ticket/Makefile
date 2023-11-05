init:
	@echo "== 👩‍🌾 init =="
	brew install go
	brew install node
	brew install pre-commit
	brew install golangci-lint
	brew upgrade golangci-lint

	@echo "== pre-commit setup =="
	pre-commit install

	@echo "== test lib setup=="
	go install github.com/Thatooine/go-test-html-report@latest
	go install github.com/vakenbolt/go-test-report@latest
	go install golang.org/x/tools/cmd/goimports@latest

test:
	@echo "== 🦸‍️ ci.tester =="
	go install github.com/vakenbolt/go-test-report@v0.9.3
	go install github.com/Thatooine/go-test-html-report@v1.1.0
	go test -v ./... -tags dynamic -coverprofile=coverage.out -json | go-test-report
	go test -v ./... -tags dynamic -coverprofile=coverage.out -json | go-test-html-report
	go tool cover -html=coverage.out -o test_report_coverage.html
	open test_report_coverage.html
	open test_report.html
	open report.html

start:
	go run --tags dynamic ./cmd/main.go serve-rest

mock:
	go generate -x ./mocks/mock_gen.go

build:
	@echo "== 🛠‍️ build =="
	go mod tidy && \
	go mod vendor && \
	go build -mod=vendor -a -installsuffix cgo -tags musl -o main ./cmd/main.go

ci.lint:
	@echo "== 🙆 ci.linter =="
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
	golangci-lint run -v ./... --build-tags=musl --fix