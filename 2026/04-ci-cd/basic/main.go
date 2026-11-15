package main

import (
	"fmt"
	"io"
	"os"
)

func main(){
    os.Exit( run(os.Args[1:],os.Stdout,os.Stderr) )
}

func run(args []string, stdout, stderr io.Writer) int {
	numbers, err := parseNumbers(args)
	if err != nil {
		fmt.Fprintln(stderr, "Ошибка:", err)
		fmt.Fprintln(stderr, "Использование: stats <число> [<число> ...]")
		return 1
	}

	stats, err := calculateStats(numbers)
	if err != nil {
		fmt.Fprintln(stderr, "Ошибка:", err)
		return 1
	}

	fmt.Fprintf(stdout, "Count:   %d\nSum:     %d\nMin:     %d\nMax:     %d\nAverage: %.2f\n",
		stats.Count, stats.Sum, stats.Min, stats.Max, stats.Average)
	return 0
}
