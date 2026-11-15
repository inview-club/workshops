package main

import (
	"bytes"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestCalculateStats(t *testing.T) {
	tests := []struct {
		name    string
		numbers []int64
		want    statistics
		wantErr bool
	}{
		{name: "positive", numbers: []int64{10, 20, 30, 40}, want: statistics{4, 100, 10, 40, 25}},
		{name: "negative", numbers: []int64{-10, -20, -30}, want: statistics{3, -60, -30, -10, -20}},
		{name: "mixed", numbers: []int64{-5, 0, 8}, want: statistics{3, 3, -5, 8, 1}},
		{name: "single", numbers: []int64{7}, want: statistics{1, 7, 7, 7, 7}},
		{name: "zero", numbers: []int64{0}, want: statistics{1, 0, 0, 0, 0}},
		{name: "fractional average", numbers: []int64{1, 2}, want: statistics{2, 3, 1, 2, 1.5}},
		{name: "empty", wantErr: true},
		{name: "positive overflow", numbers: []int64{math.MaxInt64, 1}, wantErr: true},
		{name: "negative overflow", numbers: []int64{math.MinInt64, -1}, wantErr: true},
		{name: "int64 boundaries", numbers: []int64{math.MaxInt64, math.MinInt64}, want: statistics{2, -1, math.MinInt64, math.MaxInt64, -0.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateStats(tt.numbers)
			if (err != nil) != tt.wantErr {
				t.Fatalf("calculateStats() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("calculateStats() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseNumbers(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []int64
		wantErr bool
	}{
		{name: "integers", args: []string{"10", "-20", "0", "+3"}, want: []int64{10, -20, 0, 3}},
		{name: "empty", wantErr: true},
		{name: "text", args: []string{"10", "hello"}, wantErr: true},
		{name: "decimal", args: []string{"1.5"}, wantErr: true},
		{name: "empty argument", args: []string{""}, wantErr: true},
		{name: "out of range", args: []string{"9223372036854775808"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNumbers(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseNumbers() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantOut string
		wantErr string
	}{
		{name: "success", args: []string{"10", "20", "30", "40"}, wantOut: "Count:   4\nSum:     100\nMin:     10\nMax:     40\nAverage: 25.00\n"},
		{name: "missing numbers", wantErr: "передайте хотя бы одно целое число"},
		{name: "invalid argument", args: []string{"10", "oops"}, wantErr: `аргумент "oops"`},
		{name: "overflow", args: []string{"9223372036854775807", "1"}, wantErr: "промежуточная сумма"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(tt.args, &stdout, &stderr)
			if stdout.String() != tt.wantOut {
				t.Errorf("stdout = %q, want %q", stdout.String(), tt.wantOut)
			}
			if tt.wantErr == "" {
				if code != 0 || stderr.Len() != 0 {
					t.Errorf("code = %d, stderr = %q; want success", code, stderr.String())
				}
			} else if code == 0 || !strings.Contains(stderr.String(), tt.wantErr) {
				t.Errorf("code = %d, stderr = %q; want failure containing %q", code, stderr.String(), tt.wantErr)
			}
		})
	}
}
