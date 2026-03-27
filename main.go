package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "gtr"}
	rootCmd.AddCommand(NewDefaultCommand())
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

func NewDefaultCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

			fmt.Printf("%s\n", "s")

		},
	}

	return cmd
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
