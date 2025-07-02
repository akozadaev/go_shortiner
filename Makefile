MODULE = $(shell go list -m)

# 🛠 Установка всех утилит
tools:
	go install github.com/mgechev/revive@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0

# 📦 Генерация всего, что помечено //go:generate
generate:
	go generate ./...

# ⚙️ Сборка сервера
build:
	CGO_ENABLED=0 go build -a -o go_shurtiner $(MODULE)/cmd/go_shurtiner

# 📦 Сборка для продакшена (Linux AMD64)
release:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o go_shurtiner $(MODULE)/cmd/go_shurtiner/
	zip -9 -r ./go_shurtiner.zip go_shurtiner

# 🧹 Форматирование gofmt (автоисправление)
fmt:
	gofmt -s -w .

# 🧪 Тестирование
test:
	go test -v ./...

# 🧪 Покрытие тестами
test-coverage:
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# 🧼 Быстрый линтинг (только gofmt)
lint:
	gofmt -w .

# 🧼 Полный линтинг с golangci-lint (версия 2)
lint-full:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.59.0; \
	fi
	golangci-lint run ./...

# 🧼 Автофиксы revive (не исправляет всё, но помогает)
lint-fix:
	revive -formatter stylish -fix ./...

# 📚 Генерация Swagger-документации
doc:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/go_shurtiner/main.go --pd --parseGoList=false --parseDepth=2 -o ./docs/v1 --instanceName v1
	swag init -g cmd/go_shurtiner/main.go --pd --parseGoList=false --parseDepth=2 -o ./docs/v2 --instanceName v2

# 🧪 Финальная проверка перед коммитом
check: fmt lint-full test
