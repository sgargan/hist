package main

import (
	"encoding/json"
	"fmt"

	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HistTests struct {
	suite.Suite
}

func (t *HistTests) TestReadValues() {
	values := ""
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			values = fmt.Sprintf("%s %d", values, (x*10)+y)
		}
		values = fmt.Sprintf("%s\n", values)
	}
	t.testRead(values, 100)

	t.testRead(`
	1
		2
			3
				4
	`, 4)
	t.testRead(`
		1,2,3,4,5,6,7,8,9,10,
		1 2 3 4 5 6 7 8 9 10,
		1	2	3	4	5	6	7	8	9	10,
		1,2,3,4,5,6,7,8,9,10,
	`, 40)
}

func (t *HistTests) TestJson() {
	values := make([]float64, 0)
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			values = append(values, float64(y))
		}
	}
	stats := calculateStats(values, 20)
	asJson, err := stats.Json()
	t.NoError(err)

	rt := &Stats{}
	t.NoError(json.Unmarshal([]byte(asJson), rt))
	t.Equal(stats, rt)
}

func (t *HistTests) TestHistogram() {
	values := make([]float64, 0)
	values = append(values, 0.0)
	values = append(values, 200.0)
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			values = append(values, 50.0)
		}
	}
	stats := calculateStats(values, 20)
	t.Equal(` 91.14 ┤                ╭╮
 86.58 ┤               ╭╯│
 82.03 ┤               │ │
 77.47 ┤               │ │
 72.91 ┤               │ │
 68.35 ┤               │ ╰╮
 63.80 ┤               │  │
 59.24 ┤              ╭╯  │
 54.68 ┤              │   │
 50.13 ┤              │   │
 45.57 ┤              │   │
 41.01 ┤              │   ╰╮
 36.46 ┤             ╭╯    │
 31.90 ┤             │     │
 27.34 ┤             │     │
 22.78 ┤             │     │
 18.23 ┤             │     ╰╮
 13.67 ┤            ╭╯      │
  9.11 ┤            │       │
  4.56 ┤            │       │
  0.00 ┼────────────╯       ╰──────────────────────────────────────────────────────────
     0.00                       50.25                      100.50                     200.00                     `, stats.Chart())
	t.Equal(`num: 102, bins: 20, min:  0, max: 200.00, mean: 50.98, median: 50.00, variance: 246.554, stddev: 15.702
q75: 50.00, q90: 50.00, q99: 50.00, q99.9: 200.00`, stats.String())
}

func (h *HistTests) testRead(values string, expected int) {
	read, err := readValues(strings.NewReader(values))
	h.NoError(err)
	h.Equal(expected, len(read))
}

func TestTestsSuite(t *testing.T) {
	suite.Run(t, new(HistTests))
}
