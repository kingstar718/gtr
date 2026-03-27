package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gtr [input]",
		Short: "Universal conversion tool - auto-detect input type",
		Long: `Universal conversion tool that automatically detects input type:
  - Coordinate: 113.2233,23.3131 or 113.2233|23.3131 or 113.2233 23.3131
  - HTTP: https://api.example.com or https://api.example.com '{"key":"value"}'
  - Timestamp: 10-digit (1727087511) or 13-digit (1727087511000)
  - Text: any other input for base64/url/md5 conversion`,
		Args: cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
				return nil
			}

			input := strings.Join(args, " ")
			return handleAutoConvert(input)
		},
	}

	rootCmd.AddCommand(NewTimeCommand())
	rootCmd.AddCommand(NewCoordinateCommand())
	rootCmd.AddCommand(NewHttpCommand())
	rootCmd.AddCommand(NewTextCommand())
	rootCmd.AddCommand(NewVersionCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("gtr version %s\n", Version)
		},
	}
}
