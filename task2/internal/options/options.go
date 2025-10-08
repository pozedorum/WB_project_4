// Package options содержит стуктуру FlagStruct, используется для парсинга флагов командной строки
package options

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

type FlagStruct struct {
	AFlag          *int
	BFlag          *int
	CFlag          *int
	SmallCFlag     *bool
	IFlag          *bool
	VFlag          *bool
	FFlag          *bool
	NFlag          *bool
	ConcurrentMode *bool
	Pattern        string
}

func ParseOptions() (*FlagStruct, []string) {
	var fs FlagStruct

	fs.AFlag = flag.IntP("A", "A", 0, "Print N lines after each match")
	fs.BFlag = flag.IntP("B", "B", 0, "Print N lines before each match")
	fs.CFlag = flag.IntP("C", "C", 0, "Print N lines around each match (A+B)")
	fs.SmallCFlag = flag.BoolP("c", "c", false, "Only print count of matching lines")
	fs.IFlag = flag.BoolP("i", "i", false, "Ignore case distinctions")
	fs.VFlag = flag.BoolP("v", "v", false, "Select non-matching lines")
	fs.FFlag = flag.BoolP("F", "F", false, "Interpret pattern as literal string")
	fs.NFlag = flag.BoolP("n", "n", false, "Print line numbers with output")

	ePattern := flag.StringP("e", "e", "", "Pattern to search for")

	// 	ФЛАГ ВКЛЮЧЕНИЯ РАСПРЕДЕЛЁННОЙ ВЕРСИИ УТИЛИТЫ
	fs.ConcurrentMode = flag.BoolP("Q", "Q", false, "Turn on concurrent mode")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] -e PATTERN [FILE...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s [OPTIONS] PATTERN [FILE...]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()

	if *ePattern != "" {
		fs.Pattern = *ePattern
	} else if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	} else {
		fs.Pattern = args[0]
		args = args[1:]
	}

	return &fs, args
}

func (fs *FlagStruct) PrintFlags() {
	fmt.Println("flag A -", *(fs.AFlag))
	fmt.Println("flag B -", *(fs.BFlag))
	fmt.Println("flag C -", *(fs.CFlag))
	fmt.Println("flag c -", *(fs.SmallCFlag))
	fmt.Println("flag i -", *(fs.IFlag))
	fmt.Println("flag v -", *(fs.VFlag))
	fmt.Println("flag F -", *(fs.FFlag))
	fmt.Println("flag i -", *(fs.IFlag))
}
