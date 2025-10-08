// Package grep содержит функцию Grep и необходимые для её работы вспомогательные функции
package grep

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/pozedorum/WB_project_4/task2/internal/options"
)

// Grep выполняет поиск по шаблону в текстовом потоке с учетом флагов
func Grep(input io.Reader, fs options.FlagStruct, writer io.Writer) error {

	var re *regexp.Regexp
	var err error

	if *fs.FFlag {
		// Фиксированная строка - экранируем спецсимволы
		fs.Pattern = regexp.QuoteMeta(fs.Pattern)
	}

	if *fs.IFlag {
		// Игнорирование регистра
		re, err = regexp.Compile("(?i)" + fs.Pattern)
	} else {
		re, err = regexp.Compile(fs.Pattern)
	}

	if err != nil {
		return fmt.Errorf("invalid fs.Pattern: %v", err)
	}

	// Читаем все строки в память для обработки контекста
	scanner := bufio.NewScanner(input)
	lines := make([]string, 0)
	lineNumbers := make([]int, 0)
	lineIdx := 1

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		lineNumbers = append(lineNumbers, lineIdx)
		lineIdx++
	}

	if err = scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	// Определяем, какие строки совпадают с шаблоном
	matches := make([]bool, len(lines))
	for i, line := range lines {
		matches[i] = re.MatchString(line)
		if *fs.VFlag {
			matches[i] = !matches[i]
		}
	}

	// Если флаг -c, просто считаем совпадения
	if *fs.SmallCFlag {
		count := 0
		for _, match := range matches {
			if match {
				count++
			}
		}
		fmt.Fprintf(writer, "%d\n", count)
		return nil
	}

	// Определяем контекст на основе флагов
	before := 0
	after := 0

	if *fs.CFlag > 0 {
		before = *fs.CFlag
		after = *fs.CFlag
	}
	if *fs.BFlag > 0 {
		before = *fs.BFlag
	}
	if *fs.AFlag > 0 {
		after = *fs.AFlag
	}

	// Обрабатываем контекст и выводим результат
	printed := make(map[int]bool)
	var output bytes.Buffer

	for i, isMatch := range matches {
		if isMatch {
			start := max(0, i-before)
			end := min(len(lines)-1, i+after)

			for j := start; j <= end; j++ {
				if !printed[j] {
					if *fs.NFlag {
						// Определяем тип строки: совпадение или контекст
						if matches[j] { // Это строка совпадения
							output.WriteString(fmt.Sprintf("%d:", lineNumbers[j]))
						} else { // Это контекстная строка
							output.WriteString(fmt.Sprintf("%d-", lineNumbers[j]))
						}
					}
					output.WriteString(lines[j] + "\n")
					printed[j] = true
				}
			}
		}
	}

	// Выводим результат
	_, err = writer.Write(output.Bytes())
	return err
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
