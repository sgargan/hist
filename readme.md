# Hist cli

A simple cli to generate stats and a histogram from command line input using [gonum](https://gonum.org/v1/gonum) and [asciigraph](https://github.com/guptarohit/asciigraph )

# Installation

To install run
```
go install github.com/sgargan/hist@latest
```

# Running

It's intended to be used through stdin e.g.

```
python -c 'import numpy as np; [print(x) for x in np.random.normal(50, 10, 10000)]' | ./hist
 521 ┤                                     ╭─╮
 495 ┤                                    ╭╯ ╰─╮
 469 ┤                                ╭───╯    ╰╮
 443 ┤                               ╭╯         ╰╮
 417 ┤                              ╭╯           ╰╮
 391 ┤                              │             │
 365 ┤                           ╭╮╭╯             ╰╮
 339 ┤                          ╭╯╰╯               ╰─╮
 313 ┤                          │                    ╰╮
 287 ┤                         ╭╯                     ╰╮
 261 ┤                        ╭╯                       │
 235 ┤                       ╭╯                        ╰╮
 209 ┤                      ╭╯                          ╰─╮
 182 ┤                     ╭╯                             ╰─╮
 156 ┤                    ╭╯                                ╰╮
 130 ┤                   ╭╯                                  │
 104 ┤                 ╭─╯                                   ╰─╮
  78 ┤                ╭╯                                       ╰╮
  52 ┤              ╭─╯                                         ╰──╮
  26 ┤         ╭────╯                                              ╰───╮
   0 ┼─────────╯                                                       ╰─────────────
     13.91                      32.98                      52.05                      89.18
num: 10000, min: 14, max: 89.18, mean: 50.17, median: 50.20, variance: 100.391, stddev: 10.020
q75: 56.90, q90: 63.26, q99: 73.40, q99.9: 79.56
```

By default, hist will calculate the number of buckets for the histogram using the Freedman-Diaconis rule. This can be overridden by passing the `-b` flag

Generating the chart can be skipped by pass the `-c` flag

All the generated stats can be emitted to json using `-j`, the chart is omitted when formatting with json