package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// ReadProcessWrite читает файл inputPath, передает содержимое в process, записывает результат
// в outputPath, при ошибке на любом этапе возвращает ее, обернув контекстным сообщением
func ReadProcessWrite(
	inputPath string,
	outputPath string,
	process func(string) (string, error),
) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("открытие входного файла: %w", err)
	}
	defer inputFile.Close()

	data, err := io.ReadAll(inputFile)
	if err != nil {
		return fmt.Errorf("чтение файла: %w", err)
	}

	processed, err := process(string(data))
	if err != nil {
		return fmt.Errorf("обработка текста: %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("создание выходного файла: %w", err)
	}
	defer outputFile.Close()

	if _, err := outputFile.WriteString(processed); err != nil {
		return fmt.Errorf("запись в файл: %w", err)
	}

	return nil
}

func main() {
	processFunc := func(s string) (string, error) {
		return strings.ToUpper(s), nil
	}

	if err := ReadProcessWrite("input.txt", "output.txt", processFunc); err != nil {
		log.Fatalf("Программа завершилась с ошибкой: %v", err)
	}

	fmt.Println("Файл успешно обработан и сохранён в output.txt")
}
