package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/guptarohit/asciigraph"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

// Stats defines useful statistics about a set of values which can be converted to json
type Stats struct {
	Bins      int       `json:"bins"`
	Mean      float64   `json:"mean"`
	Median    float64   `json:"median"`
	Max       float64   `json:"max"`
	Min       float64   `json:"min"`
	Variance  float64   `json:"variance"`
	Stddev    float64   `json:"stddev"`
	Q75       float64   `json:"q75"`
	Q90       float64   `json:"q90"`
	Q99       float64   `json:"q99"`
	Q999      float64   `json:"q999"`
	Xs        []float64 `json:"xs"`
	Histogram Histogram `json:"histogram"`
}

// Histogram creates a
type Histogram struct {
	Dividers []float64 `json:"dividers"`
	Values   []float64 `json:"values"`
}

// String creates a string representation of the stats with
func (s *Stats) String() string {
	return fmt.Sprintf(`num: %d, bins: %d, min: %2.f, max: %.2f, mean: %.2f, median: %.2f, variance: %.3f, stddev: %.3f
q75: %.2f, q90: %.2f, q99: %.2f, q99.9: %.2f`, len(s.Xs), s.Bins, s.Min, s.Max, s.Mean, s.Median, s.Variance, s.Stddev, s.Q75, s.Q90, s.Q99, s.Q999)
}

func (s *Stats) Json() (string, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

func (s *Stats) Chart() string {
	divs := []string{toF(s.Min), "", "", toF(s.Max)}
	quarter := len(s.Histogram.Dividers) / 4
	total := 0
	for x := 1; x < 3; x++ {
		divs[x] = toF(s.Histogram.Dividers[x*quarter])
		total += len(divs[x])
	}

	return fmt.Sprintf("%s\n     %-27s%-27s%-27s%-27s",
		asciigraph.Plot(s.Histogram.Values,
			asciigraph.Height(20),
			asciigraph.Width(80),
			asciigraph.Offset(0)),
		divs[0], divs[1], divs[2], divs[3])
}

func toF(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func calculateStats(xs []float64, bins int) *Stats {
	sort.Float64s(xs)

	min := floats.Min(xs)
	max := floats.Max(xs)
	q75 := stat.Quantile(0.75, stat.Empirical, xs, nil)

	if bins == 0 {
		// calculate the number of bins for the chart using the Freedman-Diaconis rule
		q25 := stat.Quantile(0.25, stat.Empirical, xs, nil)
		iqr := q75 - q25
		binWidth := (2 * iqr) / (math.Pow(float64(len(xs)), 0.33333))
		bins = int(math.Ceil((max - min) / binWidth))
	}

	variance := stat.Variance(xs, nil)
	dividers := make([]float64, bins+1)
	floats.Span(dividers, min, max+1)

	return &Stats{
		Histogram: Histogram{
			Dividers: dividers,
			Values:   stat.Histogram(nil, dividers, xs, nil),
		},
		Bins:     bins,
		Max:      max,
		Min:      min,
		Mean:     stat.Mean(xs, nil),
		Median:   stat.Quantile(0.5, stat.Empirical, xs, nil),
		Q75:      q75,
		Q90:      stat.Quantile(0.90, stat.Empirical, xs, nil),
		Q99:      stat.Quantile(0.99, stat.Empirical, xs, nil),
		Q999:     stat.Quantile(0.999, stat.Empirical, xs, nil),
		Stddev:   math.Sqrt(variance),
		Variance: variance,
		Xs:       xs,
	}
}

func readValues(reader io.Reader) ([]float64, error) {
	lines, err := readLines(reader)
	if err != nil {
		return nil, err
	}
	values := make([]float64, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			value, err := strconv.ParseFloat(line, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing value %s: %s", line, err)
			}
			values = append(values, value)
		}
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("no values found via stdin")
	}
	return values, nil
}

func readLines(reader io.Reader) (lines []string, err error) {
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading from stdin - %v", err)
	}

	values := []string{string(all)}
	for _, delim := range []string{"\n", " ", ",", "\t"} {
		values = splitAll(delim, values...)
	}
	return values, nil
}

func splitAll(delim string, s ...string) []string {
	values := make([]string, 0)

	for _, v := range s {
		values = append(values, strings.Split(v, delim)...)
	}
	return values
}
