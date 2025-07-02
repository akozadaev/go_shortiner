// Package credential содержит тестовые утилиты для проверки системы авторизации.
package credential

// BasicAuthHeaderValue пароль для тестов
const BasicAuthHeaderValue = "Basic YWtvemFkYWV2QGluYm94LnJ1OlBAJCR3MHJk"

func getTestCredentials() string {
	return "akozadaev@inbox.ru:P@$$w0rd"
}
