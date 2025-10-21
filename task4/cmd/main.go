package main

import "github.com/pozedorum/WB_project_4/task4/pkg/analyzer"

func main() {
	analyzer := analyzer.NewAnalyzer()
	analyzer.PrintMemStats()
}
