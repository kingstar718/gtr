package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func isTimeFormat(input string) bool {
	trimmed := strings.TrimSpace(input)
	if trimmed != input {
		return false
	}

	length := len(input)

	if length == 10 || length == 13 {
		if _, err := strconv.ParseInt(input, 10, 64); err == nil {
			return true
		}
	}

	timeFormats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"20060102150405",
		"20060102150400",
		"20060102",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
		"01/02/2006 15:04:05",
		"01/02/2006 15:04",
		"01/02/2006",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		time.RFC3339,
		time.RFC822,
	}

	for _, format := range timeFormats {
		if _, err := time.Parse(format, input); err == nil {
			return true
		}
	}

	return false
}

func handleTimeConvert(input string) error {
	originArg := strings.TrimSpace(input)

	timeFormats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"20060102150405",
		"20060102150400",
		"20060102",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
		"01/02/2006 15:04:05",
		"01/02/2006 15:04",
		"01/02/2006",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		time.RFC3339,
		time.RFC822,
	}

	var parseTime time.Time
	parsed := false
	l := len(originArg)

	if !parsed && l == 10 {
		if parseInt, err := strconv.ParseInt(originArg, 10, 64); err == nil {
			parseTime = time.Unix(parseInt, 0)
			parsed = true
		}
	}

	if !parsed && l == 13 {
		if parseInt, err := strconv.ParseInt(originArg, 10, 64); err == nil {
			parseTime = time.UnixMilli(parseInt)
			parsed = true
		}
	}

	if !parsed {
		for _, format := range timeFormats {
			if t, err := time.Parse(format, originArg); err == nil {
				parseTime = t
				parsed = true
				break
			}
		}
	}

	if !parsed {
		parseTime = time.Now()
	}

	timestamp13 := parseTime.UnixNano() / int64(time.Millisecond)
	timestamp10 := parseTime.Unix()

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "                input\t:\t%s\n", originArg)
	fmt.Fprintf(w, "          timestamp10\t:\t%d\n", timestamp10)
	fmt.Fprintf(w, "          timestamp13\t:\t%d\n", timestamp13)
	fmt.Fprintf(w, "  2006-01-02 15:04:05\t:\t%s\n", parseTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "       20060102150405\t:\t%s\n", parseTime.Format("20060102150405"))
	fmt.Fprintf(w, "  2006/01/02 15:04:05\t:\t%s\n", parseTime.Format("2006/01/02 15:04:05"))
	fmt.Fprintf(w, "              RFC3339\t:\t%s\n", parseTime.Format(time.RFC3339))

	w.Flush()
	fmt.Print(buf.String())

	return nil
}

func NewTimeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "time",
		Aliases: []string{"t", "-t", "--t", "-time", "--time"},
		Short: "\n-------------------------------\n" +
			"| COMMAND: time               |\n" +
			"| TYPE: Time Format Convert   |\n" +
			"| INPUT:                      |\n" +
			"|   1. 10 timestamp           |\n" +
			"|   2. 13 timestamp           |\n" +
			"|   3. standard date          |\n" +
			"| EXAMPLES:                   |\n" +
			"|   1. 1727087511             |\n" +
			"|   2. 1727087511000          |\n" +
			"|   3. 2024-09-23 10:31:51    |\n" +
			"|   4. 20240923103151         |\n" +
			"-------------------------------\n",
		RunE: func(cmd *cobra.Command, args []string) error {
			var originArg string
			if len(args) == 1 {
				originArg = args[0]
			} else if len(args) == 2 {
				originArg = args[0] + " " + args[1]
			}

			return handleTimeConvert(originArg)
		},
	}
	return cmd
}

