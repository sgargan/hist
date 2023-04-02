package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type HistCommand struct {
	*cobra.Command
	chart bool
	json  bool
	bins  int
}

func NewHistCommand() *HistCommand {
	hist := &HistCommand{}
	cmd := &cobra.Command{
		Use:   "hist",
		Short: "a simple cli for generating histograms",
		Long:  `a simple cli for generating histograms`,
		Run:   hist.Run,
	}
	hist.Command = cmd

	hist.PersistentFlags().BoolVarP(&hist.chart, "chart", "c", true, "show the histogram chart")
	hist.PersistentFlags().BoolVarP(&hist.json, "json", "j", false, "output stats details as json, does not render chart")
	hist.PersistentFlags().IntVarP(&hist.bins, "bins", "b", 0, "no of bins to use for the histogram. Calculates automatically if not specified")
	hist.Command = cmd

	return hist
}

func (h *HistCommand) Run(cmd *cobra.Command, args []string) {
	xs, err := readValues(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading values: %s", err)
		os.Exit(1)
	}

	stats := calculateStats(xs, h.bins)
	if h.json {
		asJson, err := stats.Json()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating json: %s", err)
			os.Exit(1)
		}
		fmt.Println(asJson)
	} else if h.chart {
		fmt.Println(stats.Chart())
	}
	fmt.Println(stats.String())
}

func main() {
	NewHistCommand().Execute()
}
