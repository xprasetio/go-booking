.PHONY: test test-coverage test-package test-user test-user-service open-coverage install-mockery generate-mocks

# Install mockery
install-mockery:
	@echo "Installing mockery..."
	@go install github.com/vektra/mockery/v2@latest

# Generate mocks
generate-mocks:
	@echo "Generating mocks..."
	@mockery --dir=./internal/user --name=UserServiceInterface --output=./internal/user/mocks --outpkg=mocks
	@mockery --dir=./pkg/logger --name=Logger --output=./internal/user/mocks --outpkg=mocks
	@mockery --dir=./pkg/redis --name=RedisClient --output=./internal/user/mocks --outpkg=mocks

# Menjalankan semua test
test:
	go test -v ./...

# Menjalankan test dengan coverage dan generate HTML report
test-coverage:
	@mkdir -p tmp
	go test -coverprofile=tmp/coverage.out ./...
	go tool cover -html=tmp/coverage.out -o tmp/coverage.html
	@echo "Coverage report generated at tmp/coverage.html"
	@open tmp/coverage.html

# Menjalankan test untuk package tertentu dengan coverage
test-package:
	@if [ -z "$(package)" ]; then \
		echo "Usage: make test-package package=<package_path>"; \
		exit 1; \
	fi
	@mkdir -p tmp
	go test -v -coverprofile=tmp/coverage.out $(package)
	go tool cover -html=tmp/coverage.out -o tmp/coverage.html
	@echo "Coverage report generated at tmp/coverage.html"
	@open tmp/coverage.html

# Menjalankan test untuk domain user dengan coverage
test-user:
	@mkdir -p tmp
	go test -v -coverprofile=tmp/coverage.out ./internal/user/...
	go tool cover -html=tmp/coverage.out -o tmp/coverage.html
	@echo "Coverage report generated at tmp/coverage.html"
	@open tmp/coverage.html

# Menjalankan test khusus untuk user_service dengan coverage
test-user-service:
	@mkdir -p tmp
	go test -v -coverprofile=tmp/coverage.out ./internal/user/user_service_test.go
	go tool cover -html=tmp/coverage.out -o tmp/coverage.html
	@echo "Coverage report generated at tmp/coverage.html"
	@open tmp/coverage.html

# Membuka file coverage.html di browser
open-coverage:
	@if [ ! -f "tmp/coverage.html" ]; then \
		echo "Coverage report not found. Please run test-coverage first."; \
		exit 1; \
	fi
	@open tmp/coverage.html 