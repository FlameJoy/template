package initializers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnv загружает переменные окружения из файла
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Разделяем строку по символу '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in .env file: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Устанавливаем переменную окружения
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting environment variable: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}
