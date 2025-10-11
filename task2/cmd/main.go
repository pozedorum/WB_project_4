package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pozedorum/WB_project_4/task2/internal/concurrency"
	"github.com/pozedorum/WB_project_4/task2/internal/grep"
	"github.com/pozedorum/WB_project_4/task2/internal/options"
)

func main() {
	fs, fileArgs := options.ParseOptions()
	if len(fileArgs) <= 0 {
		fmt.Printf("grep: no files detected")
		return
	}
	if *fs.ConcurrentMode > 0 {
		// Распределённый режим
		var files []*os.File

		if len(fileArgs) == 0 {
			// Используем stdin
			files = append(files, os.Stdin)
		} else {
			// Открываем файлы
			for _, filename := range fileArgs {
				file, err := os.Open(filename)
				if err != nil {
					// log.Printf("Warning: cannot open %s: %v", filename, err)
					continue
				}
				defer func() {
					if err := file.Close(); err != nil {
						fmt.Printf("grep: %v", err)
					}
				}()
				files = append(files, file)
			}
		}

		if len(files) == 0 {
			log.Fatal("No files to process")
		}

		// Создаем мастера (например, с 4 воркерами)
		master, err := concurrency.NewMaster(*fs.ConcurrentMode, fs)
		if err != nil {
			log.Fatal(err)
		}

		// Обрабатываем файлы
		if err := master.ProcessFilesStreaming(files, "grep", fs.Pattern); err != nil {
			log.Fatal(err)
		}

		// Выводим результаты
		results := master.MergeResults()
		for _, line := range results {
			fmt.Println(line)
		}

	} else {
		// Обработка файла (аргумент - имя файла)
		for _, fileName := range fileArgs {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatalf("grep: %s: %v", file.Name(), err)
			}
			defer func() {
				if err := file.Close(); err != nil {
					fmt.Printf("grep: %v", err)
				}
			}()

			err = grep.Grep(file, *fs, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
