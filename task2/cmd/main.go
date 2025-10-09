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
		fmt.Errorf("grep: no files detected")
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
					log.Printf("Warning: cannot open %s: %v", filename, err)
					continue
				}
				defer file.Close()
				files = append(files, file)
			}
		}

		if len(files) == 0 {
			log.Fatal("No files to process")
		}

		// Создаем мастера (например, с 4 воркерами)
		master, err := concurrency.NewMaster(4, fs)
		if err != nil {
			log.Fatal(err)
		}
		defer master.Close()

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
		for _, file := range fileArgs {
			file, err := os.Open(file)
			if err != nil {
				log.Fatalf("grep: %s: %v", file.Name(), err)
			}
			defer file.Close()

			err = grep.Grep(file, *fs, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
