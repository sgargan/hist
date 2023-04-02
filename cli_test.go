package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CliTests struct {
	suite.Suite
}

func (t *CliTests) TestJson() {
	cli := NewHistCommand()
	cli.ParseFlags([]string{"hist", "-j"})
	t.Equal(true, cli.json)
}

func (t *CliTests) TestChart() {
	cli := NewHistCommand()
	cli.ParseFlags([]string{"hist", "-c"})
	t.Equal(true, cli.chart)
}

func (t *CliTests) TestBins() {
	cli := NewHistCommand()
	cli.ParseFlags([]string{"hist", "-b", "20"})
	t.Equal(20, cli.bins)
}

func TestCliTestsSuite(t *testing.T) {
	suite.Run(t, new(CliTests))
}
