package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type statistics struct {
	Count int
	Sum     int64
	Min     int64
	Max     int64
	Average float64
}

func parseNumbers(args []string) ([]int64, error) {
	if len(args) == 0 {
		return nil, errors.New("передайте хотя бы одно целое число")
	}

	numbers := make([]int64, len(args))
	for i, arg := range args {
		number, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("аргумент %q должен быть целым числом в диапазоне int64", arg)
		}
		numbers[i] = number
	}
	return numbers, nil
}

func calculateStats(numbers []int64) (statistics, error) {
	if len(numbers) == 0 {
		return statistics{}, errors.New("набор чисел не должен быть пустым")
	}

	stats := statistics{Count: len(numbers), Min: numbers[0], Max: numbers[0]}
	for _, number := range numbers {
		// Проверяем переполнение до сложения, чтобы не получить неверную сумму.
		if (number > 0 && stats.Sum > math.MaxInt64-number) ||
			(number < 0 && stats.Sum < math.MinInt64-number) {
			return statistics{}, errors.New("промежуточная сумма выходит за диапазон int64")
		}
		stats.Sum+=number
		if number < stats.Min {
			stats.Min = number
		}
		if number > stats.Max {
			stats.Max = number
		}
	}
	stats.Average=float64(stats.Sum)/float64(stats.Count)
	return stats, nil
}
